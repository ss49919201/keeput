{-# LANGUAGE DeriveGeneric #-}
{-# LANGUAGE OverloadedStrings #-}

module Service.Notifier where

import Data.Aeson
import Data.Aeson.Types (parseMaybe)
import GHC.Generics (Generic)
import Network.HTTP.Simple
import System.Environment (lookupEnv)

newtype ReqBody = ReqBody {text :: String}
  deriving (Generic)

instance ToJSON ReqBody

getIsGoalAchieved :: Object -> Maybe Bool
getIsGoalAchieved = parseMaybe $ \obj -> obj .: "is_goal_achieved"

reqBody :: Bool -> ReqBody
reqBody True = ReqBody "目標達成です🎊よく頑張りました！"
reqBody False = ReqBody "目標未達です😢これから頑張りましょう！"

loadValueFromEnv :: String -> IO (Maybe String)
loadValueFromEnv key = do
  maybeUrl <- lookupEnv key
  return $ case maybeUrl of
    Just url | not (null url) -> Just url
    _ -> Nothing

loadSlackWebhookUrl :: IO (Maybe String)
loadSlackWebhookUrl = loadValueFromEnv "SLACK_WEBHOOK_URL"

notificationCondition :: IO (Maybe String)
notificationCondition = do
  loadValueFromEnv "NOTIFICATION_CONDITION"

data NotificationCondition
  = OnFailure
  | OnSuccess
  | Always

parseNotificationCondition :: Maybe String -> NotificationCondition
parseNotificationCondition (Just "on_failure") = OnFailure
parseNotificationCondition (Just "on_success") = OnFailure
parseNotificationCondition _ = Always

sendReq :: ReqBody -> IO (Either String ())
sendReq reqBody = do
  maybeUrl <- loadSlackWebhookUrl
  case maybeUrl of
    Nothing -> do
      return $ Left "slack webhook url is not set or empty"
    Just url -> do
      request <- parseRequest $ "POST " ++ url
      let requestWithBody = setRequestBodyJSON reqBody request
      response <- httpLBS requestWithBody
      let status = getResponseStatusCode response
      print status
      if status == 200
        then do
          print "success request"
          return $ Right ()
        else do
          let errorMsg = "failed to request with status code: " ++ show status
          putStrLn errorMsg
          putStrLn $ "Response body: " ++ show (getResponseBody response)
          return $ Left errorMsg

notify :: Object -> IO ()
notify value = do
  case getIsGoalAchieved value of
    Just isGoalAchieved -> do
      result <- sendReq (reqBody isGoalAchieved)
      either (putStrLn . ("sendReq failed: " ++)) (const $ return ()) result
    Nothing -> putStrLn "failed to get isGoalAchieved"

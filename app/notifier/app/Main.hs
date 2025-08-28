{-# LANGUAGE DeriveGeneric #-}
{-# LANGUAGE ImportQualifiedPost #-}
{-# LANGUAGE OverloadedStrings #-}

module Main where

import Data.Aeson
import Data.Aeson.Types (parseMaybe)
import Data.ByteString.Lazy qualified as LBS
import Usecase.Notifier

getJsonFromStdio :: IO (Maybe Object)
getJsonFromStdio = decode <$> LBS.getContents

main :: IO ()
main = do
  decoded <- getJsonFromStdio
  case decoded of
    Just value -> notify value
    Nothing -> putStrLn "failed to decode json string."

import { vValidator } from "@hono/valibot-validator";
import { Hono } from "hono";
import { HTTPException } from "hono/http-exception";
import { logger } from "hono/logger";
import { object, string } from "valibot";
import { Locker } from "./locker";

type Bindings = {
  LOCKER: DurableObjectNamespace<Locker>;
  API_KEY: string;
};

const app = new Hono<{ Bindings: Bindings }>();

app.use(logger());

app.use("*", async (c, next) => {
  if (c.env.API_KEY === "") {
    console.error("env api key is empty");
    throw new HTTPException(500);
  }

  const apiKey = c.req.header("X-LOCKER-API-KEY");
  if (apiKey !== c.env.API_KEY) {
    throw new HTTPException(401);
  }
  await next();
});

app.get("health", (c) => {
  return c.json("ok");
});

app.post(
  "/acquire",
  vValidator(
    "json",
    object({
      lockId: string(),
    })
  ),
  async (c) => {
    const { lockId } = c.req.valid("json");
    const stub = c.env.LOCKER.get(c.env.LOCKER.idFromName(lockId));
    const result = await stub.acquire();
    return c.json({
      msg: result ? "ok" : "ng",
    });
  }
);

app.post(
  "/release",
  vValidator(
    "json",
    object({
      lockId: string(),
    })
  ),
  async (c) => {
    const { lockId } = c.req.valid("json");
    const stub = c.env.LOCKER.get(c.env.LOCKER.idFromName(lockId));
    await stub.release();
    return c.json({
      msg: "ok",
    });
  }
);

export default app;
export { Locker };

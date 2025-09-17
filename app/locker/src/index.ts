import { vValidator } from "@hono/valibot-validator";
import { Hono } from "hono";
import { logger } from "hono/logger";
import { object, string } from "valibot";
import { Locker } from "./locker";

type Bindings = {
  LOCKER: DurableObjectNamespace<Locker>;
};

const app = new Hono<{ Bindings: Bindings }>();

app.use(logger());

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

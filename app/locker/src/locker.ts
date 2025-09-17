import { DurableObject } from "cloudflare:workers";

// 10秒間リクエストがない場合、インスタンスが破棄される。
// この挙動を利用して一定時間後の自動リリースを実現している。
export class Locker extends DurableObject {
  #isLocked = false;

  public constructor(state: DurableObjectState, env: unknown) {
    super(state, env);
  }

  public async acquire(): Promise<boolean> {
    if (this.#isLocked) {
      return false;
    }

    this.#isLocked = true;
    return true;
  }

  public async release(): Promise<void> {
    this.#isLocked = false;
  }
}

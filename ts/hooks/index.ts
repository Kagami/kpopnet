import { EventEmitter } from "events";

export const enum HOOKS {
  showAlert,
}

const hooks = new EventEmitter();

export function hook(name: HOOKS, fn: (...args: any[]) => void) {
  hooks.addListener(name.toString(), fn);
}

export function unhook(name: HOOKS, fn: (...args: any[]) => void) {
  hooks.removeListener(name.toString(), fn);
}

export function trigger(name: HOOKS, ...args: any[]) {
  hooks.emit(name.toString(), ...args);
}

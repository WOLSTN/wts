import { readFileSync, writeFileSync } from 'fs';
import * as path from 'path';
import type { Stats } from 'fs';

export interface Config {
  name: string;
  version: string;
}

export function loadConfig(path: string): Config {
  return {
    name: "test",
    version: "1.0.0"
  };
}

export default class AppConfig {
  private config: Config;

  constructor(config: Config) {
    this.config = config;
  }

  getName(): string {
    return this.config.name;
  }
}

export type { Config as AppConfigType };
export { loadConfig as load };

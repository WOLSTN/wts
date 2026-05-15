enum Color {
  Red,
  Green,
  Blue
}

enum Status {
  Active = 1,
  Inactive = 0,
  Pending = 2
}

const enum Direction {
  Up = "UP",
  Down = "DOWN",
  Left = "LEFT",
  Right = "RIGHT"
}

type Point = {
  x: number;
  y: number;
};

type Vector2D = Point & { z: number };

type ID = string | number;

type Callback<T> = (value: T) => void;

namespace Utils {
  export function log(msg: string): void {
    console.log(msg);
  }

  export const VERSION = "1.0.0";

  export interface Config {
    debug: boolean;
  }
}

namespace Math {
  export function square(x: number): number {
    return x * x;
  }
}

const c: Color = Color.Red;
const status: Status = Status.Active;
const dir: Direction = Direction.Up;
const point: Point = { x: 1, y: 2 };
const id: ID = "test-123";

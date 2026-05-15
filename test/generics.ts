interface Container<T> {
  value: T;
  getValue(): T;
}

interface Pair<T, U> {
  first: T;
  second: U;
}

function identity<T>(arg: T): T {
  return arg;
}

function makePair<T, U>(first: T, second: U): Pair<T, U> {
  return { first, second };
}

class Box<T> {
  private value: T;

  constructor(value: T) {
    this.value = value;
  }

  getValue(): T {
    return this.value;
  }
}

const numContainer: Container<number> = {
  value: 42,
  getValue() { return this.value; }
};

const strContainer: Container<string> = {
  value: "hello",
  getValue() { return this.value; }
};

const num = identity(42);
const str = identity("hello");
const pair = makePair(1, "two");
const box = new Box<number>(100);

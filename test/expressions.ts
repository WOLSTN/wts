const add = (a: number, b: number): number => a + b;

const multiply = (x: number, y: number) => {
  return x * y;
};

const result = add(1, 2);

const max = (a: number, b: number) => a > b ? a : b;

const min = a < b ? a : b;

const arr1 = [1, 2, 3];
const arr2 = [...arr1, 4, 5];

const obj1 = { a: 1, b: 2 };
const obj2 = { ...obj1, c: 3 };

const typeStr = typeof 42;

async function asyncFunc() {
  const value = await Promise.resolve(42);
  return value;
}

function* generatorFunc() {
  yield 1;
  yield 2;
}

const nonNull = obj1!.a;

const asserted = value as string;

const deleted = delete obj1.a;

const voided = void 0;

class Foo {
  constructor() {
    console.log(new.target);
  }
}

const tagged = tag`Hello ${name}!`;

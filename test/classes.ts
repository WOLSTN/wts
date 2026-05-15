interface IPoint {
  x: number;
  y: number;
  distance(): number;
}

interface IShape {
  name: string;
  area(): number;
}

class Point implements IPoint {
  public x: number;
  public y: number;

  constructor(x: number, y: number) {
    this.x = x;
    this.y = y;
  }

  distance(): number {
    return Math.sqrt(this.x * this.x + this.y * this.y);
  }
}

class Rectangle extends Point implements IShape {
  private width: number;
  protected height: number;
  readonly name: string = "Rectangle";

  constructor(x: number, y: number, width: number, height: number) {
    super(x, y);
    this.width = width;
    this.height = height;
  }

  area(): number {
    return this.width * this.height;
  }
}

abstract class Shape {
  abstract getArea(): number;
  
  describe(): string {
    return "A shape";
  }
}

const p = new Point(1, 2);
const r = new Rectangle(0, 0, 10, 20);

interface Point {
    x: number;
    y: number;
}

function add(a: number, b: number): number {
    return a + b;
}

function greet(name: string): string {
    return `Hello, ${name}!`;
}

const point: Point = { x: 1, y: 2 };

class Calculator {
    private value: number = 0;
    
    add(n: number): void {
        this.value += n;
    }
    
    getValue(): number {
        return this.value;
    }
}

const calc = new Calculator();
calc.add(10);

export { add, greet, Calculator, Point };

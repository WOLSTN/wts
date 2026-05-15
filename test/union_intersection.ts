type StringOrNumber = string | number;
type NameAndAge = { name: string } & { age: number };

function formatValue(value: StringOrNumber): string {
    return String(value);
}

function printPerson(person: NameAndAge): void {
    console.log(`${person.name}: ${person.age}`);
}

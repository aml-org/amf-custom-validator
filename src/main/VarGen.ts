let counter = 0;

export const genvar = (hint: string) => {
    counter++;
    return `gen_${hint}_${counter}`;
}


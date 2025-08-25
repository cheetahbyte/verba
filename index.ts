const path = "./main.verba";
const file = Bun.file(path);

const arrbuff = await file.arrayBuffer();
const buffer = Buffer.from(arrbuff);
console.log(buffer.toString());

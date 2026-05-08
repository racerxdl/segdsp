export const waterfallLut: [number, number, number][] = [];

for (let i = 0; i < 256; i++) {
  let r: number, g: number, b: number;
  if (i < 20) {
    r = 0; g = 0; b = 0;
  } else if (i < 70) {
    r = 0; g = 0; b = Math.round(140 * (i - 20) / 50);
  } else if (i < 100) {
    r = Math.round(60 * (i - 70) / 30);
    g = Math.round(125 * (i - 70) / 30);
    b = Math.round(115 * (i - 70) / 30 + 140);
  } else if (i < 150) {
    r = Math.round(195 * (i - 100) / 50 + 60);
    g = Math.round(130 * (i - 100) / 50 + 125);
    b = Math.round(255 - 255 * (i - 100) / 50);
  } else if (i < 250) {
    r = 255;
    g = Math.round(255 - 255 * (i - 150) / 100);
    b = 0;
  } else {
    r = 255;
    g = Math.round(255 * (i - 250) / 5);
    b = Math.round(255 * (i - 250) / 5);
  }
  waterfallLut.push([r, g, b]);
}

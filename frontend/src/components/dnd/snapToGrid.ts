export function snapToGrid(x: number, y: number): [number, number] {
  const snappedX = Math.round(x);
  const snappedY = Math.round(y);
  return [snappedX, snappedY];
}

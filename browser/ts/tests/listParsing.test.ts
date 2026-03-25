import { describe, expect, it } from "vitest";
import { parseListItemsFromResponse } from "../src/crud/listParsing";

describe("parseListItemsFromResponse", () => {
  it("returns array as-is", () => {
    expect(parseListItemsFromResponse([{ id: "a" }])).toEqual([{ id: "a" }]);
  });

  it("reads items from envelope", () => {
    expect(parseListItemsFromResponse({ items: [{ id: "x" }] })).toEqual([{ id: "x" }]);
  });

  it("treats null items as empty", () => {
    const envelope: { items: null; total: number } = { items: null, total: 0 };
    expect(parseListItemsFromResponse(envelope)).toEqual([]);
  });

  it("treats missing items as empty", () => {
    expect(parseListItemsFromResponse({})).toEqual([]);
  });
});

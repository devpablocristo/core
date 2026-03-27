type ListEnvelope<T> = {
  items?: T[] | null | unknown;
  data?: unknown;
};

/**
 * Normaliza respuestas de listado tipo arreglo, `{ items }` o envelopes BFF anidados con `data`.
 * Go suele serializar slices nil como `null`; cualquier rama sin arreglo termina en lista vacía.
 */
export function parseListItemsFromResponse<T>(input: unknown): T[] {
  const queue: unknown[] = [input];
  const seen = new Set<unknown>();

  while (queue.length > 0) {
    const current = queue.shift();
    if (Array.isArray(current)) {
      return current as T[];
    }
    if (current == null || typeof current !== "object" || seen.has(current)) {
      continue;
    }

    seen.add(current);
    const envelope = current as ListEnvelope<T>;

    if (Array.isArray(envelope.items)) {
      return envelope.items;
    }

    if ("data" in envelope) {
      queue.push(envelope.data);
    }
    if ("items" in envelope) {
      queue.push(envelope.items);
    }
  }

  return [];
}

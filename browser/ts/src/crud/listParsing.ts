/**
 * Normaliza respuestas de listado tipo { items } o arreglo plano.
 * Go suele serializar slice nil como JSON null — se trata como lista vacía.
 */
export function parseListItemsFromResponse<T>(data: { items?: T[] | null } | T[]): T[] {
  if (Array.isArray(data)) {
    return data;
  }
  const items = data.items;
  if (items == null) {
    return [];
  }
  return items;
}

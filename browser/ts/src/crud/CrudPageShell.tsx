import type { ReactNode } from "react";

/**
 * Layout canónico de consola: cabecera (título + acciones), error, formulario,
 * barra de herramientas (búsqueda / filtros) y cuerpo (spinner, vacío o tabla).
 * Productos (Pymes u otros) montan datos y i18n encima; las clases CSS viven en el tema de la app.
 * El host puede combinar `data-theme` y, si aplica, `data-admin-skin` en `<html>` para tokens adicionales.
 */
export type CrudPageShellProps = {
  title: ReactNode;
  subtitle?: ReactNode;
  /** Bajo el título, columna izquierda (p. ej. filtros tipo píldora). */
  headerLeadSlot?: ReactNode;
  /** Columna derecha de cabecera (p. ej. búsqueda + fila de botones); el tema aplica `.crud-page-shell__header-actions`. */
  headerActions?: ReactNode;
  error?: ReactNode;
  /** Tarjeta de formulario alta/edición */
  form?: ReactNode;
  /** Fila búsqueda + archivados u otros filtros */
  toolbar?: ReactNode;
  children: ReactNode;
};

export function CrudPageShell({
  title,
  subtitle,
  headerLeadSlot,
  headerActions,
  error,
  form,
  toolbar,
  children,
}: CrudPageShellProps) {
  return (
    <>
      <div className="page-header crud-page-shell__header">
        <div className="crud-page-shell__header-main">
          <h1 className="crud-page-shell__title">{title}</h1>
          {subtitle != null && subtitle !== false ? (
            <p className="text-secondary">{subtitle}</p>
          ) : null}
          {headerLeadSlot != null ? headerLeadSlot : null}
        </div>
        {headerActions != null ? (
          <div className="crud-page-shell__header-actions">{headerActions}</div>
        ) : null}
      </div>
      {error}
      {form}
      {toolbar}
      {children}
    </>
  );
}

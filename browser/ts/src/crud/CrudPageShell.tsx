import type { ReactNode } from "react";

/**
 * Layout canónico de consola: cabecera (título + acciones), error, formulario,
 * barra de herramientas (búsqueda / filtros) y cuerpo (spinner, vacío o tabla).
 * Productos (Pymes u otros) montan datos y i18n encima; las clases CSS viven en el tema de la app.
 */
export type CrudPageShellProps = {
  title: ReactNode;
  subtitle?: ReactNode;
  /** Botones a la derecha del título (CSV, crear, etc.) */
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
  headerActions,
  error,
  form,
  toolbar,
  children,
}: CrudPageShellProps) {
  return (
    <>
      <div className="page-header">
        <div>
          <h1>{title}</h1>
          {subtitle != null && subtitle !== false ? (
            <p className="text-secondary">{subtitle}</p>
          ) : null}
        </div>
        {headerActions != null ? (
          <div className="actions-row">{headerActions}</div>
        ) : null}
      </div>
      {error}
      {form}
      {toolbar}
      {children}
    </>
  );
}

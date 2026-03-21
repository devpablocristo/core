export interface TokenPair {
  access_token: string;
  refresh_token?: string | null;
}

export interface BrowserTokenStorageOptions {
  namespace: string;
  storage?: Storage;
  hostAware?: boolean;
  accessTokenKey?: string;
  refreshTokenKey?: string;
  legacyKeys?: string[];
}

export interface BrowserTokenStorage {
  key(name: string): string;
  getAccessToken(): string | null;
  getRefreshToken(): string | null;
  setAccessToken(token: string): void;
  setRefreshToken(token: string): void;
  setTokens(tokens: TokenPair): void;
  clear(): void;
}

function resolveStorage(custom?: Storage): Storage | null {
  if (custom) {
    return custom;
  }
  if (typeof window === "undefined") {
    return null;
  }
  return window.localStorage;
}

function resolvePrefix(namespace: string, hostAware: boolean): string {
  const cleanNamespace = namespace.trim();
  if (!hostAware || typeof window === "undefined") {
    return `${cleanNamespace}:`;
  }
  return `${cleanNamespace}:${window.location.host}:`;
}

export function createBrowserTokenStorage(options: BrowserTokenStorageOptions): BrowserTokenStorage {
  const storage = resolveStorage(options.storage);
  const hostAware = options.hostAware ?? true;
  const accessTokenKey = options.accessTokenKey ?? "access_token";
  const refreshTokenKey = options.refreshTokenKey ?? "refresh_token";
  const legacyKeys = options.legacyKeys ?? [accessTokenKey, refreshTokenKey];

  function key(name: string): string {
    return `${resolvePrefix(options.namespace, hostAware)}${name}`;
  }

  function get(name: string): string | null {
    return storage?.getItem(key(name)) ?? null;
  }

  function set(name: string, value: string): void {
    storage?.setItem(key(name), value);
  }

  function clear(): void {
    if (!storage) {
      return;
    }

    legacyKeys.forEach((legacyKey) => storage.removeItem(legacyKey));

    const prefix = resolvePrefix(options.namespace, hostAware);
    const toRemove: string[] = [];
    for (let index = 0; index < storage.length; index += 1) {
      const itemKey = storage.key(index);
      if (itemKey && itemKey.startsWith(prefix)) {
        toRemove.push(itemKey);
      }
    }
    toRemove.forEach((itemKey) => storage.removeItem(itemKey));
  }

  return {
    key,
    getAccessToken: () => get(accessTokenKey),
    getRefreshToken: () => get(refreshTokenKey),
    setAccessToken: (token: string) => set(accessTokenKey, token),
    setRefreshToken: (token: string) => set(refreshTokenKey, token),
    setTokens: (tokens: TokenPair) => {
      set(accessTokenKey, tokens.access_token);
      if (tokens.refresh_token) {
        set(refreshTokenKey, tokens.refresh_token);
      }
    },
    clear,
  };
}

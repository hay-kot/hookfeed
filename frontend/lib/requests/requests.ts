export enum Method {
  GET = "GET",
  POST = "POST",
  PUT = "PUT",
  DELETE = "DELETE",
  PATCH = "PATCH",
}

export type ResponseInterceptor = (r: Response, rq?: RequestInit) => void;
export type RequestInterceptor = (request: RequestInit) => RequestInit;

export interface TResponse<T> {
  status: number;
  error: boolean;
  data: T;
  response: Response;
}

export type RequestArgs<T> = {
  url: string;
  body?: T;
  data?: FormData;
  headers?: Record<string, string>;
};

export class Requests {
  private headers: Record<string, string> = {};
  private responseInterceptors: ResponseInterceptor[] = [];
  private requestInterceptors: RequestInterceptor[] = [];
  private bearer: () => string | null = () => null;

  withResponseInterceptor(interceptor: ResponseInterceptor) {
    this.responseInterceptors.push(interceptor);
    return this;
  }

  withRequestInterceptor(interceptor: RequestInterceptor) {
    this.requestInterceptors.push(interceptor);
    return this;
  }

  withBearer(bearer: () => string | null) {
    this.bearer = bearer;
    return this;
  }

  private callResponseInterceptors(response: Response, request?: RequestInit) {
    this.responseInterceptors.forEach((i) => i(response, request));
  }

  private applyRequestInterceptors(request: RequestInit): RequestInit {
    return this.requestInterceptors.reduce(
      (modifiedRequest, interceptor) => interceptor(modifiedRequest),
      request,
    );
  }

  constructor(headers: Record<string, string> = {}) {
    this.headers = headers;
  }

  public get<T>(args: RequestArgs<T>): Promise<TResponse<T>> {
    return this.do<T>(Method.GET, args);
  }

  public post<T, U>(args: RequestArgs<T>): Promise<TResponse<U>> {
    return this.do<U>(Method.POST, args);
  }

  public put<T, U>(args: RequestArgs<T>): Promise<TResponse<U>> {
    return this.do<U>(Method.PUT, args);
  }

  public delete<T>(args: RequestArgs<T>): Promise<TResponse<T>> {
    return this.do<T>(Method.DELETE, args);
  }

  public patch<T, U>(args: RequestArgs<T>): Promise<TResponse<U>> {
    return this.do<U>(Method.PATCH, args);
  }

  public async download(url: string, filename: string) {
    // fetches the url and downloads the contents as a file
    try {
      const initialRequest: RequestInit = {
        method: "GET",
        headers: {
          Authorization: `Bearer ${this.bearer()}`,
          "Content-Type": "application/octet-stream", // Adjust if needed
        },
      };

      const modifiedRequest = this.applyRequestInterceptors(initialRequest);

      const response = await fetch(url, modifiedRequest);
      if (!response.ok) {
        throw new Error(`Network response was not ok: ${response.statusText}`);
      }
      const blob = await response.blob();
      const link = document.createElement("a");
      link.href = URL.createObjectURL(blob);
      link.download = filename;
      link.click();
      // Clean up the object URL after the download
      URL.revokeObjectURL(link.href);
    } catch (e) {
      console.error(e);
    }
  }

  private methodSupportsBody(method: Method): boolean {
    return (
      method === Method.POST || method === Method.PUT || method === Method.PATCH
    );
  }

  private async do<T>(
    method: Method,
    rargs: RequestArgs<unknown>,
  ): Promise<TResponse<T>> {
    const payload: RequestInit = {
      method,
      headers: {
        ...rargs.headers,
        ...this.headers,
      } as Record<string, string>,
    };

    if (this.methodSupportsBody(method)) {
      if (rargs.data) {
        payload.body = rargs.data;
      } else {
        // @ts-expect-error - we know that the header is there
        payload.headers["Content-Type"] = "application/json";
        payload.body = JSON.stringify(rargs.body);
      }
    }

    const maybeBearer = this.bearer();
    if (maybeBearer) {
      // @ts-expect-error - we know that the header is there
      payload.headers.Authorization = `Bearer ${maybeBearer}`;
    }

    // Apply request interceptors
    const modifiedPayload = this.applyRequestInterceptors(payload);

    const response = await fetch(rargs.url, modifiedPayload);
    this.callResponseInterceptors(response, modifiedPayload);

    const data: T = await (async () => {
      if (response.status === 204) {
        return {} as T;
      }
      if (
        response.headers.get("Content-Type")?.startsWith("application/json")
      ) {
        try {
          return await response.json();
        } catch (e) {
          return {} as T;
        }
      }
      return response.body as unknown as T;
    })();

    return {
      status: response.status,
      error: !response.ok,
      data,
      response,
    };
  }
}

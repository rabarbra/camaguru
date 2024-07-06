import { API_ENDPOINT } from "../config.js";

/**
 * Represents an API client for making HTTP requests
 * to the backend server pong.
 * @class
 */
class ApiClient {
  constructor (baseUrl) {
    this.baseUrl = new URL(baseUrl);
    this.route = null;
    this.headers = {
      Accept: 'application/json',
    };
    if (localStorage.getItem('access_token')) {
      this.headers['Authorization'] = `Bearer ${localStorage.getItem('access_token')}`;
    }
    if (localStorage.getItem('refresh_token')) {
      this.refresh_token = localStorage.getItem('refresh_token');
    }
  }

  async proceedResponse(response) {
    let data = await response.json();
    if (response.ok)
      return data;
    else
    {
      if (data.err)
      return {err: data.err};
      else if (data.status)
        return {err: data.status};
      return {err: response.status}
    }
  }

  async sendRequest (path, method, body, query, headers = {}) {
    // if (!path.startsWith('/'))
    //   path = '/' + path;
    const url = new URL(this.baseUrl + path);
    console.log(url)
    let params = {
      method,
      headers: {...this.headers, ...headers},
    };
    if (body) {
      params.body = JSON.stringify(body);
    }
    params.headers['Content-Type'] = 'application/json';
    // params.mode = "no-cors"
    if (query) {
      Object.entries(query).forEach(([key, value]) => {
        url.searchParams.append(key, value);
      });
    }
    try {
      let response = await fetch(url, params);
      return await this.proceedResponse(response);
    } catch (error) {
      return {err: error.message};
    }
  }

  async get (path, query) {
    return await this.sendRequest(path, 'GET', null, query);
  }

  async post (path, body, query, headers) {
    return await this.sendRequest(path, 'POST', body, query, headers);
  }

  async put (path, body, query) {
    return await this.sendRequest(path, 'PUT', body, query);
  }

  async delete (path, body, query) {
    return await this.sendRequest(path, 'DELETE', body, query);
  }

  async authorize (payload, query = null) {
    let response_body;
    if ("access_token" in payload)
      response_body = payload;
    else {
      response_body = await this.sendRequest('signin', 'POST', payload, query);
    }
    if (response_body.error)
      return response_body;
    const access_token = response_body.access_token;
    this.headers['Authorization'] = `Bearer ${access_token}`;
    localStorage.setItem('access_token', access_token);
    const me = await this.get("/me");
    if (me.error)
      return me;
    localStorage.setItem("me", JSON.stringify(me));
    return {"ok": "true"};
  }

  unauthorize () {
    localStorage.removeItem('access_token');
    localStorage.removeItem('me');
  }

  authorized() {
    return (
      (
        localStorage.getItem('access_token')
      ) && localStorage.getItem('me')
      ? true
      : false
    );
  }

}

export const apiClient = new ApiClient(API_ENDPOINT);

export default ApiClient;
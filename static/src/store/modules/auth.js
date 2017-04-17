import jwtDecode from 'jwt-decode';
import axios from 'axios';
import * as types from '../mutation-types';
import auth from '../../api/auth';

const jwtStorageKey = 'token';

const decodeToken = function(token) {
  let data = null;

  if (token) {
    try {
      data = jwtDecode(token);
    } catch (e) {
      return null;
    }
  }

  return data;
};

const store = {
  state: {
    authUrl: null,
    token: localStorage.getItem(jwtStorageKey),
  },
  getters: {
    isLoggedIn({ token }) {
      const data = decodeToken(token);

      if (data && data.github_token) {
        return true;
      }

      return false;
    },
    name({ token }) {
      const data = decodeToken(token);
      return data ? data.name : null;
    },
    avatarUrl({ token }) {
      const data = decodeToken(token);
      return data ? data.avatar_url : null;
    },
    authRedirectUrl: ({ authUrl }) => authUrl,
  },
  actions: {
    requestAuth({ commit }) {
      auth
        .requestAuth()
        .then(({ data }) => {
          const { url, token } = data;
          commit(types.REQUEST_AUTH_SUCCESS, { url, token });
        })
        .catch(() => commit(types.REQUEST_AUTH_FAILURE));
    },

    processAuth({ commit }, { state, code }) {
      auth
        .processAuth(state, code)
        .then(({ data }) => {
          const { token } = data;
          commit(types.PROCESS_AUTH_SUCCESS, { token });
        })
        .catch(() => commit(types.PROCESS_AUTH_FAILURE));
    },

    logout({ commit }) {
      commit(types.CLEAR_AUTH_TOKEN);
    },
  },
  mutations: {
    [types.REQUEST_AUTH_SUCCESS](state, { url, token }) {
      state.authUrl = url;
      state.token = token;
      localStorage.setItem(jwtStorageKey, token);
    },

    [types.PROCESS_AUTH_SUCCESS](state, { token }) {
      state.token = token;
      localStorage.setItem(jwtStorageKey, token);
    },

    [types.CLEAR_AUTH_TOKEN](state) {
      state.token = null;
      localStorage.setItem(jwtStorageKey, null);
    },
  },
};

axios.interceptors.request.use(config => {
  const { token } = store.state;

  if (token) {
    const authorisedConfig = config;
    authorisedConfig.headers.Authorization = `Bearer ${token}`;
    return authorisedConfig;
  }

  return config;
});

export default store;

import jwtDecode from 'jwt-decode';
import axios from 'axios';
import * as types from '../mutation-types';
import auth from '../../api/auth';

const jwtStorageKey = 'token';

const decodeToken = function(token) {
  if (!token) {
    return null;
  }

  try {
    return jwtDecode(token);
  } catch (e) {
    return null;
  }
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
    async requestAuth({ commit }) {
      try {
        const response = await auth.requestAuth();
        const { url, token } = response.data;
        commit(types.REQUEST_AUTH_SUCCESS, {
          url,
          token,
        });
      } catch (error) {
        commit(types.REQUEST_AUTH_FAILURE, { error });
      }
    },

    async processAuth({ commit }, { state, code }) {
      try {
        const response = await auth.processAuth(state, code);
        const { token } = response.data;
        commit(types.PROCESS_AUTH_SUCCESS, { token });
      } catch (error) {
        commit(types.PROCESS_AUTH_FAILURE, { error });
      }
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

import * as types from '../mutation-types';
import exports from '../../api/exports';

export default {
  state: {
    exports: [],
    loadingExports: false,
  },

  getters: {
    exports: state => state.exports,
    loadingExports: state => state.loadingExports,
  },

  actions: {
    async listExports({ commit }, { titleId, rplId }) {
      commit(types.LIST_EXPORTS_LOADING);

      try {
        const response = await exports.getExports(titleId, rplId);
        commit(types.LIST_EXPORTS_SUCCESS, {
          exports: response.data,
        });
      } catch (error) {
        commit(types.LIST_EXPORTS_FAILURE, { error });
      }
    },

    clearExports({ commit }) {
      commit(types.CLEAR_EXPORTS);
    },
  },

  mutations: {
    [types.LIST_EXPORTS_LOADING](state) {
      state.loadingExports = true;
    },

    [types.LIST_EXPORTS_SUCCESS](state, { exports }) {
      state.exports = exports;
      state.loadingExports = false;
    },

    [types.LIST_EXPORTS_FAILED](state) {
      state.loadingExports = false;
    },

    [types.CLEAR_EXPORTS](state) {
      state.exports = [];
      state.loadingExports = false;
    },
  },
};

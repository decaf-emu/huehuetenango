import * as types from '../mutation-types';
import imports from '../../api/imports';

export default {
  state: {
    imports: [],
    loadingImports: false,
  },

  getters: {
    imports: state => state.imports,
    loadingImports: state => state.loadingImports,
  },

  actions: {
    async listImports({ commit }, { titleId, rplId }) {
      commit(types.LIST_IMPORTS_LOADING);

      try {
        const response = await imports.getImports(titleId, rplId);
        commit(types.LIST_IMPORTS_SUCCESS, {
          imports: response.data,
        });
      } catch (error) {
        commit(types.LIST_IMPORTS_FAILURE, { error });
      }
    },

    clearImports({ commit }) {
      commit(types.CLEAR_IMPORTS);
    },
  },

  mutations: {
    [types.LIST_IMPORTS_LOADING](state) {
      state.loadingImports = true;
    },

    [types.LIST_IMPORTS_SUCCESS](state, { imports }) {
      state.imports = imports;
      state.loadingImports = false;
    },

    [types.LIST_IMPORTS_FAILED](state) {
      state.loadingImports = false;
    },

    [types.CLEAR_IMPORTS](state) {
      state.imports = [];
      state.loadingImports = false;
    },
  },
};

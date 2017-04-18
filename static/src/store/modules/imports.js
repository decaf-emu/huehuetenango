import * as types from '../mutation-types';
import imports from '../../api/imports';

const state = {
  imports: [],
  loadingImports: false,
};

const getters = {
  imports: state => state.imports,
  loadingImports: state => state.loadingImports,
};

const actions = {
  listImports({ commit }, { titleId, rplId }) {
    commit(types.LIST_IMPORTS_LOADING);

    imports
      .getImports(titleId, rplId)
      .then(({ data }) => commit(types.LIST_IMPORTS_SUCCESS, { imports: data }))
      .catch(() => commit(types.LIST_IMPORTS_FAILURE));
  },

  clearImports({ commit }) {
    commit(types.CLEAR_IMPORTS);
  },
};

const mutations = {
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
  },
};

export default {
  state,
  getters,
  actions,
  mutations,
};

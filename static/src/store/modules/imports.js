import * as types from '../mutation-types';
import imports from '../../api/imports';

const state = {
  imports: [],
};

const getters = {
  imports: state => state.imports,
};

const actions = {
  listExports({ commit }, { titleId, rplId }) {
    commit(types.LIST_IMPORTS_LOADING);

    imports.getImports(
      titleId,
      rplId,
      results => commit(types.LIST_IMPORTS_SUCCESS, { results }),
      () => commit(types.LIST_IMPORTS_FAILURE),
    );
  },
};

const mutations = {
  [types.LIST_IMPORTS_LOADING](state) {
    state.imports = [];
  },
  [types.LIST_IMPORTS_SUCCESS](state, { results }) {
    state.imports = results;
  },
};

export default {
  state,
  getters,
  actions,
  mutations,
};

import * as types from '../mutation-types';
import imports from '../../api/imports';

const state = {
  imports: [],
};

const getters = {
  imports: state => state.imports,
};

const actions = {
  listImports({ commit }, { titleId, rplId }) {
    commit(types.LIST_IMPORTS_LOADING);

    imports
      .getImports(titleId, rplId)
      .then(({ data }) => commit(types.LIST_IMPORTS_SUCCESS, { imports: data }))
      .catch(() => commit(types.LIST_IMPORTS_FAILURE));
  },
};

const mutations = {
  [types.LIST_IMPORTS_LOADING](state) {
    state.imports = [];
  },
  [types.LIST_IMPORTS_SUCCESS](state, { imports }) {
    state.imports = imports;
  },
};

export default {
  state,
  getters,
  actions,
  mutations,
};

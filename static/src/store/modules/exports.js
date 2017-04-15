import * as types from '../mutation-types';
import exports from '../../api/exports';

const state = {
  exports: [],
};

const getters = {
  exports: state => state.exports,
};

const actions = {
  listExports({ commit }, { titleId, rplId }) {
    commit(types.LIST_EXPORTS_LOADING);

    exports.getExports(
      titleId,
      rplId,
      results => commit(types.LIST_EXPORTS_SUCCESS, { results }),
      () => commit(types.LIST_EXPORTS_FAILURE),
    );
  },
};

const mutations = {
  [types.LIST_EXPORTS_LOADING](state) {
    state.exports = [];
  },
  [types.LIST_EXPORTS_SUCCESS](state, { results }) {
    state.exports = results;
  },
};

export default {
  state,
  getters,
  actions,
  mutations,
};

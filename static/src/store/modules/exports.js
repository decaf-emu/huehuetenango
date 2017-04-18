import * as types from '../mutation-types';
import exports from '../../api/exports';

const state = {
  exports: [],
  loadingExports: false,
};

const getters = {
  exports: state => state.exports,
  loadingExports: state => state.loadingExports,
};

const actions = {
  listExports({ commit }, { titleId, rplId }) {
    commit(types.LIST_EXPORTS_LOADING);

    exports
      .getExports(titleId, rplId)
      .then(({ data }) => commit(types.LIST_EXPORTS_SUCCESS, { exports: data }))
      .catch(() => commit(types.LIST_EXPORTS_FAILURE));
  },

  clearExports({ commit }) {
    commit(types.CLEAR_EXPORTS);
  },
};

const mutations = {
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
};

export default {
  state,
  getters,
  actions,
  mutations,
};

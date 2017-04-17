import * as types from '../mutation-types';
import rpls from '../../api/rpls';

const state = {
  titleRpls: [],
  rpl: null,
};

const getters = {
  titleRpls: state => state.titleRpls,
  rpl: state => state.rpl,
};

const actions = {
  getTitleRpls({ commit }, titleId) {
    commit(types.TITLE_RPLS_LOADING);

    rpls
      .getTitleRpls(titleId)
      .then(({ data }) => commit(types.TITLE_RPLS_SUCCESS, { rpls: data }))
      .catch(() => commit(types.TITLE_RPLS_FAILURE));
  },
  getRpl({ commit }, { titleId, rplId }) {
    commit(types.GET_RPL_LOADING);

    rpls
      .getRpl(titleId, rplId)
      .then(({ data }) => commit(types.GET_RPL_SUCCESS, { rpl: data }))
      .catch(() => commit(types.GET_RPL_FAILURE));
  },
};

const mutations = {
  [types.TITLE_RPLS_LOADING](state) {
    state.titleRpls = [];
  },
  [types.TITLE_RPLS_SUCCESS](state, { rpls }) {
    state.titleRpls = rpls.sort((a, b) => {
      if (a.IsRPX) {
        return -1;
      }

      const nameA = a.Name.toUpperCase();
      const nameB = b.Name.toUpperCase();

      if (nameA < nameB) {
        return -1;
      }
      if (nameA > nameB) {
        return 1;
      }

      return 0;
    });
  },
  [types.GET_RPL_LOADING](state) {
    state.rpl = null;
  },
  [types.GET_RPL_SUCCESS](state, { rpl }) {
    state.rpl = rpl;
  },
};

export default {
  state,
  getters,
  actions,
  mutations,
};

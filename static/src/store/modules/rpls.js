import * as types from '../mutation-types';
import rpls from '../../api/rpls';

export default {
  state: {
    titleRpls: [],
    loadingTitleRpls: false,

    rpl: null,
    loadingRpl: false,
  },

  getters: {
    titleRpls: state => state.titleRpls,
    loadingTitleRpls: state => state.loadingTitleRpls,
    rpl: state => state.rpl,
    loadingRpl: state => state.loadingRpl,
  },

  actions: {
    async getTitleRpls({ commit }, titleId) {
      commit(types.TITLE_RPLS_LOADING);

      try {
        const response = await rpls.getTitleRpls(titleId);
        commit(types.TITLE_RPLS_SUCCESS, {
          rpls: response.data,
        });
      } catch (error) {
        commit(types.TITLE_RPLS_FAILURE, { error });
      }
    },

    clearTitleRpls({ commit }) {
      commit(types.CLEAR_TITLE_RPLS);
    },

    async getRpl({ commit }, { titleId, rplId }) {
      commit(types.GET_RPL_LOADING);

      try {
        const response = await rpls.getRpl(titleId, rplId);
        commit(types.GET_RPL_SUCCESS, {
          rpl: response.data,
        });
      } catch (error) {
        commit(types.GET_RPL_FAILURE, { error });
      }
    },

    clearRpl({ commit }) {
      commit(types.CLEAR_RPL);
    },
  },

  mutations: {
    [types.TITLE_RPLS_LOADING](state) {
      state.loadingTitleRpls = true;
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

      state.loadingTitleRpls = false;
    },

    [types.TITLE_RPLS_FAILED](state) {
      state.loadingTitleRpls = false;
    },

    [types.CLEAR_TITLE_RPLS](state) {
      state.titleRpls = [];
      state.loadingTitleRpls = false;
    },

    [types.GET_RPL_LOADING](state) {
      state.loadingRpl = true;
    },

    [types.GET_RPL_SUCCESS](state, { rpl }) {
      state.rpl = rpl;
      state.loadingRpl = false;
    },

    [types.GET_RPL_FAILED](state) {
      state.loadingRpl = false;
    },

    [types.CLEAR_RPL](state) {
      state.rpl = null;
      state.loadingRpl = false;
    },
  },
};

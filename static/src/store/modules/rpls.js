import * as types from '../mutation-types';
import rpls from '../../api/rpls';
import imports from '../../api/imports';
import exports from '../../api/exports';

export default {
  state: {
    titleRpls: [],
    loadingTitleRpls: false,

    rpl: null,
    imports: [],
    exports: [],
    loadingRpl: false,
  },

  getters: {
    titleRpls: state => state.titleRpls,
    loadingTitleRpls: state => state.loadingTitleRpls,

    rpl: state => state.rpl,
    imports: state => state.imports,
    exports: state => state.exports,
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
        const [
          rplResponse,
          importsResponse,
          exportsResponse,
        ] = await Promise.all([
          rpls.getRpl(titleId, rplId),
          imports.getImports(titleId, rplId),
          exports.getExports(titleId, rplId),
        ]);

        commit(types.GET_RPL_SUCCESS, {
          rpl: rplResponse.data,
          imports: importsResponse.data,
          exports: exportsResponse.data,
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

    [types.GET_RPL_SUCCESS](state, { rpl, imports, exports }) {
      state.rpl = rpl;
      state.imports = imports;
      state.exports = exports;
      state.loadingRpl = false;
    },

    [types.GET_RPL_FAILED](state) {
      state.loadingRpl = false;
    },

    [types.CLEAR_RPL](state) {
      state.rpl = null;
      state.imports = [];
      state.exports = [];
      state.loadingRpl = false;
    },
  },
};

import * as types from '../mutation-types';
import rpls from '../../api/rpls';
import imports from '../../api/imports';
import exports from '../../api/exports';
import sortBy from 'lodash-es/sortBy';

export default {
  state: {
    titleRpls: [],
    loadingTitleRpls: false,

    rpl: null,
    importData: [],
    importFunctions: [],
    exportData: [],
    exportFunctions: [],
    loadingRpl: false,
  },

  getters: {
    titleRpls: state => state.titleRpls,
    loadingTitleRpls: state => state.loadingTitleRpls,

    rpl: state => state.rpl,
    importData: state => state.importData,
    importFunctions: state => state.importFunctions,
    exportData: state => state.exportData,
    exportFunctions: state => state.exportFunctions,
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
      state.titleRpls = sortBy(rpls, title => title.Name.toUpperCase());
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

      if (imports && imports.data) {
        state.importData = sortBy(imports.data || [], item => item.name.toUpperCase());
      } else {
        state.importData = [];
      }

      if (imports && imports.functions) {
        state.importFunctions = sortBy(imports.functions || [], item => item.name.toUpperCase());
      } else {
        state.importFunctions = [];
      }

      if (exports && exports.data) {
        state.exportData = sortBy(exports.data || [], item => item.toUpperCase());
      } else {
        state.exportData = [];
      }

      if (exports && exports.functions) {
        state.exportFunctions = sortBy(exports.functions || [], item => item.toUpperCase());
      } else {
        state.exportFunctions = [];
      }

      state.loadingRpl = false;
    },

    [types.GET_RPL_FAILED](state) {
      state.loadingRpl = false;
    },

    [types.CLEAR_RPL](state) {
      state.rpl = null;
      state.importData = [];
      state.importFunctions = [];
      state.exportData = [];
      state.exportFunctions = [];
      state.loadingRpl = false;
    },
  },
};

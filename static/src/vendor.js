import Vue from 'vue';
import VueRouter from 'vue-router';
import VueHead from 'vue-head';
import Vuikit from 'vuikit';
import UIkit from 'uikit';
import Icons from 'uikit/dist/js/uikit-icons';

Vue.use(VueRouter);
Vue.use(VueHead, {
  separator: '-',
  complement: 'Huehuetenango',
});
Vue.use(Vuikit);
UIkit.use(Icons);

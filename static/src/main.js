import Vue from 'vue';
import VueRouter from 'vue-router';
import store from './store';
import App from './components/App.vue';
import TitleList from './components/titles/TitleList.vue';
import TitleView from './components/titles/TitleView.vue';
import Import from './components/import/Import.vue';
import './main.scss';

const router = new VueRouter({
  mode: 'history',
  routes: [
    { path: '/', name: 'home', redirect: '/titles' },
    { path: '/titles', name: 'titles', component: TitleList },
    {
      path: '/titles/:titleId/:rplId?/:type?',
      name: 'title',
      component: TitleView,
      props: true,
    },
    { path: '/import', name: 'import', component: Import },
  ],
});

new Vue({
  router,
  store,
  el: '#app',
  render: h => h(App),
});

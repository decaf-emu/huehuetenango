<template>
  <div class="uk-container uk-container-expand uk-position-relative">
    <div class="uk-margin uk-margin-top uk-margin-bottom" uk-grid>
      <h1 v-if="title">{{ title.LongNameEnglish }}</h1>
    </div>

    <div uk-grid>
      <div class="uk-width-1-5">
        <RplList class="rpl-list" v-if="titleId" :titleId="titleId" :rpls="titleRpls" :type="type" />
        <div v-if="loadingTitleRpls" class="uk-overlay uk-overlay-default uk-position-cover uk-position-z-index">
          <div class="uk-position-center" uk-spinner></div>
        </div>
      </div>
      <div class="uk-width-4-5">
        <RplView v-if="titleId && rplId" :titleId="titleId" :rplId="rplId" :type="type" />
      </div>
    </div>

    <div v-if="loadingTitle" class="uk-overlay uk-overlay-default uk-position-cover uk-position-z-index">
      <div class="uk-position-center" uk-spinner></div>
    </div>
  </div>
</template>

<script>
import { mapGetters } from 'vuex';
import RplList from '../rpls/RplList.vue';
import RplView from '../rpls/RplView.vue';

export default {
  components: { RplList, RplView },
  props: ['titleId', 'rplId', 'type'],
  beforeMount() {
    const { titleId } = this;

    if (titleId) {
      this.$store.dispatch('getTitle', titleId);
      this.$store.dispatch('getTitleRpls', titleId);
    }
  },
  computed: {
    ...mapGetters(['title', 'loadingTitle', 'titleRpls', 'loadingTitleRpls']),
  },
  watch: {
    titleId(titleId) {
      if (titleId) {
        this.$store.dispatch('getTitle', titleId);
        this.$store.dispatch('getTitleRpls', titleId);
      }
    },
    titleRpls(rpls) {
      const { titleId, rplId, type } = this;

      if (titleId && !rplId && rpls && rpls.length > 0) {
        this.$router.push({
          name: 'title',
          params: { titleId, rplId: rpls[0].ID, type },
        });
      }
    },
  },
};
</script>

<style>
.search-rpls {
  width: 100%;
}

.rpl-list {
  margin-top: 0;
}
</style>

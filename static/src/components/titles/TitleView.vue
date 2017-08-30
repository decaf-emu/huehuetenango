<template>
  <div>
    <div uk-height-viewport class="uk-panel uk-panel-scrollable uk-width-1-4@m uk-position-fixed uk-padding-small">
      <h3 v-if="title">{{ title.LongNameEnglish }}</h3>
      <RplList v-if="titleId" :loading="loadingTitleRpls" :titleId="titleId" :rpls="titleRpls" :type="type" />
    </div>
    <div uk-grid>
      <div class="uk-width-1-4@m"></div>
      <div class="uk-width-3-4@m">
        <RplView v-if="titleId && rplId" :titleId="titleId" :rplId="rplId" :type="type" />
      </div>
    </div>
  </div>
</template>

<script>
import { mapGetters, mapActions } from 'vuex';
import RplList from '../rpls/RplList.vue';
import RplView from '../rpls/RplView.vue';

export default {
  components: { RplList, RplView },
  props: ['titleId', 'rplId', 'type'],

  computed: {
    ...mapGetters([
      'title',
      'rpl',
      'loadingTitle',
      'titleRpls',
      'loadingTitleRpls',
    ]),
  },

  methods: {
    ...mapActions(['getTitle', 'getTitleRpls']),

    checkRpl() {
      const { titleId, titleRpls } = this;
      let { rplId } = this;

      if (titleId && !rplId && titleRpls && titleRpls.length > 0) {
        rplId = titleRpls[0].ID;

        this.$router.replace({
          name: 'title',
          params: { titleId, rplId, type: 'imports' },
        });
      }
    },

    async fetchTitle() {
      const { titleId } = this;

      if (titleId) {
        await Promise.all([this.getTitle(titleId), this.getTitleRpls(titleId)]);
      }
    },
  },

  watch: {
    titleId() {
      this.fetchTitle();
    },

    titleRpls() {
      this.checkRpl();
    },

    title() {
      this.$emit('updateHead');
    },

    rpl() {
      this.$emit('updateHead');
    },
  },

  head: {
    title() {
      let value;
      const { title, rpl } = this;

      if (title) {
        const { ShortNameEnglish, LongNameEnglish } = title;

        value = {
          inner: ShortNameEnglish ? ShortNameEnglish : LongNameEnglish,
        };

        if (rpl && rpl.TitleID == title.ID) {
          value.inner = `${rpl.Name} - ${value.inner}`;
        }
      }

      return value;
    },
  },

  beforeMount() {
    this.fetchTitle();
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

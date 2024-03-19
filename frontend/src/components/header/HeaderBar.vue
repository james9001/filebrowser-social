<template>
  <header>
    <img v-if="showLogo !== undefined" :src="logoURL" />
    <action
      v-if="showMenu !== undefined"
      class="menu-button"
      icon="menu"
      :label="$t('buttons.toggleSidebar')"
      @action="openSidebar()"
    />

    <slot />

    <div id="dropdown" :class="{ active: this.currentPromptName === 'more' }">
      <slot name="actions" />
    </div>

    <action
      v-if="this.$slots.actions"
      id="more"
      icon="more_vert"
      :label="$t('buttons.more')"
      @action="openMore()"
    />

    <div
      class="overlay"
      v-show="this.currentPromptName == 'more'"
      @click="$store.commit('closeHovers')"
    />
  </header>
</template>

<script>
import { logoURL } from "@/utils/constants";

import Action from "@/components/header/Action.vue";
import { mapGetters } from "vuex";

export default {
  name: "header-bar",
  props: ["showLogo", "showMenu", "hideNotifications"],
  components: {
    Action,
  },
  data: function () {
    return {
      logoURL,
    };
  },
  methods: {
    openSidebar() {
      this.$emit("update:hideNotifications");
      this.$store.commit("showHover", "sidebar");
    },
    openMore() {
      this.$emit("update:hideNotifications");
      this.$store.commit("showHover", "more");
    },
  },
  computed: {
    ...mapGetters(["currentPromptName"]),
  },
};
</script>

<style></style>

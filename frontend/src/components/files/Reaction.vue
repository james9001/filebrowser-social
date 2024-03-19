<template>
  <div class="reaction">
    <div class="existing-reactions reaction-grid">
      <div
        v-for="existingReaction in thisCommentExistingReactions"
        :key="existingReaction.id"
        class="reaction-item"
      >
        <img
          :src="getReactionImageSrc(existingReaction.reactionValue)"
          :title="getTooltipForExistingReaction(existingReaction)"
        />
      </div>
      <div
        class="reaction-item add-new-reaction noselect"
        @click="toggleAddingMode($event)"
      >
        <i class="material-icons">add_reaction</i>
      </div>
    </div>
    <div class="reaction-picker reaction-grid" v-if="inAddingMode">
      <div class="reacting-hint">React...</div>
      <div
        v-for="availableReaction in reactionsAvailable"
        :key="availableReaction"
        class="reaction-item"
        v-bind:class="{
          darkenedReaction: getExistingReactionOfSameKind(availableReaction),
        }"
        @click="addOrRemoveMyReaction(availableReaction)"
      >
        <img
          :src="getReactionImageSrc(availableReaction)"
          :title="availableReaction"
        />
      </div>
    </div>
  </div>
</template>

<script>
import { files as filesApi } from "@/api";
export default {
  name: "reaction-element",
  props: ["reactionsAvailable", "commentId", "existingReactions", "myUsername"],
  computed: {
    thisCommentExistingReactions() {
      return this.existingReactions
        ? this.existingReactions.filter(
            (reaction) => reaction.commentId === this.commentId
          )
        : [];
    },
  },
  async mounted() {
    //placeholder
  },
  beforeDestroy() {
    //placeholder
  },
  data: function () {
    return {
      inAddingMode: false,
    };
  },
  methods: {
    async toggleAddingMode() {
      this.inAddingMode = !this.inAddingMode;
    },
    getReactionImageSrc(availableReaction) {
      const imagePath = encodeURI(
        `/filebrowser-social/reactions/${availableReaction}.png`
      );
      return filesApi.getDownloadURL({ path: imagePath }, "thumb");
    },
    addOrRemoveMyReaction(availableReaction) {
      this.toggleAddingMode();
      const existingSameReaction =
        this.getExistingReactionOfSameKind(availableReaction);
      if (existingSameReaction) {
        this.$emit("update:deleteReaction", existingSameReaction.id);
      } else {
        const newReaction = {
          reactionValue: availableReaction,
          commentId: this.commentId,
        };
        this.$emit("update:newReaction", newReaction);
      }
    },
    getExistingReactionOfSameKind(availableReaction) {
      return this.existingReactions
        ? this.existingReactions.find(
            (reaction) =>
              reaction.reactionValue === availableReaction &&
              reaction.userName === this.myUsername &&
              reaction.commentId === this.commentId
          )
        : null;
    },
    getTooltipForExistingReaction(existingReaction) {
      return `${existingReaction.userName} reacted with ${existingReaction.reactionValue}`;
    },
  },
};
</script>
<style>
.reaction-grid {
  min-height: 2.75em;
  padding: 0.25em;
  background: var(--surfaceSecondary);
  border-radius: 0.5em;
  margin: 0.25em;

  display: grid;
  grid-template-columns: repeat(auto-fill, 2.25em);
}

.reaction-item {
  height: 2em;
  margin: 0.125em;

  flex: 1;
}

.reaction-item img {
  float: left;
  width: 2em;
  height: 2em;
  background-color: black;
}

.add-new-reaction {
  text-align: center;
}

.add-new-reaction i {
  color: var(--textSecondary);
  font-size: 2rem !important;
}

.reacting-hint {
  color: var(--textSecondary);
  grid-column: 1 / -1;
  font-style: italic;
  font-size: 0.75em;
}

.darkenedReaction {
  border: 0.125em solid red;
}

.darkenedReaction img {
  filter: grayscale(1);
  width: 1.75em;
  height: 1.75em;
}
</style>

<template>
  <div class="comment">
    <div class="commentUserPicture">
      <img :src="getCommentUserPictureImgSrc()" />
    </div>
    <div class="commentContent">
      <div class="commentUserAndText">
        <div class="commentUser">{{ userName }}</div>
        <div class="commentText">{{ commentText }}</div>
      </div>
      <div class="commentFunctions">
        <a
          v-if="userName === myUsername"
          class="commentDeleteButton"
          @click="deleteComment($event, id)"
        >
          Delete
        </a>
        <div class="commentTime" :title="formatIso8601AbsoluteTime(createdTime)">
          {{ formatFromNowTime(createdTime) }}
        </div>
      </div>
    </div>
    <reaction-element
      v-bind:reactionsAvailable="reactionsAvailable"
      v-bind:comment-id="id"
      v-bind:existing-reactions="existingReactions"
      v-bind:myUsername="myUsername"
      @update:newReaction="createNewReaction"
      @update:deleteReaction="deleteExistingReaction"
    >
    </reaction-element>
  </div>
</template>

<script>
import { files as filesApi } from "@/api";
import moment from "moment";
import ReactionElement from "@/components/files/Reaction.vue";

export default {
  name: "comment-element",
  components: {
    ReactionElement,
  },
  props: [
    "id",
    "filePath",
    "commentText",
    "userName",
    "createdTime",
    "myUsername",
    "reactionsAvailable",
    "existingReactions",
  ],
  async mounted() {
    //placeholder
  },
  beforeDestroy() {
    //placeholder
  },
  data: function () {
    //placeholder
    return {};
  },
  methods: {
    formatFromNowTime: function (timeValue) {
      return moment(timeValue).fromNow();
    },
    formatIso8601AbsoluteTime: function (timeValue) {
      return moment(timeValue).format("YYYY-MM-DD HH:mm");
    },
    async deleteComment(event, deletingId) {
      event.preventDefault();
      this.$emit("update:deleteComment", deletingId);
    },
    getCommentUserPictureImgSrc() {
      const imagePath = encodeURI(
        `/filebrowser-social/users/${this.userName}.png`
      );
      return filesApi.getDownloadURL({ path: imagePath }, "thumb");
    },
    async createNewReaction(newReaction) {
      this.$emit("update:newReaction", newReaction);
    },
    async deleteExistingReaction(deletingId) {
      this.$emit("update:deleteReaction", deletingId);
    },
  },
};
</script>
<style>
.comment {
  margin: 0em 0.5em 0.5em 0.5em;
  padding: 0.5em;
  background: var(--surfacePrimary);
  border-radius: 1em;
  min-height: 5em;
}

.commentUserPicture img {
  float: left;
  width: 4em;
  height: 4em;
  background-color: black;
  border-radius: 3em;
}

.commentContent {
  padding-left: 0.5em;
  min-height: 4em;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.commentUser {
  color: var(--textSecondary);
}

.commentText {
  white-space: pre-line;
}

.commentFunctions {
  font-size: 0.8em;
}

.commentDeleteButton {
  cursor: pointer;
  color: var(--blue);
}

.commentTime {
  color: var(--textSecondary);
}

.commentEntering {
  margin: 0em 0.5em 0.5em 0.5em;
  padding: 1em;
  background: var(--surfacePrimary);
  border-radius: 1em;
}

.commentEnteringButtonDiv {
  text-align: right;
}

.commentEnteringHeader {
  color: var(--textSecondary);
  font-size: 0.8em;
  padding-bottom: 0.5em;
}

#commentEnteringTextBox {
  line-height: 1.15;
  min-height: 10em;
  resize: vertical;
}
</style>

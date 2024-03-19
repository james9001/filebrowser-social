<template>
  <div class="social">
    <div
      v-if="!showComments && !isFullDesktopMode"
      class="showCommentsPrompt noselect"
    >
      <button class="button button--flat" @click="toggleComments($event, true)">
        Show {{ commentsResponse.comments.length }} comments
      </button>
    </div>
    <div v-if="showComments || isFullDesktopMode" class="socialSection">
      <div class="hideCommentsPrompt noselect" v-if="!isFullDesktopMode">
        <button
          class="button button--flat"
          @click="toggleComments($event, false)"
        >
          Hide comments
        </button>
      </div>
      <div class="inner-social-section">
        <reaction-element
          v-bind:reactionsAvailable="reactionsAvailable"
          v-bind:comment-id="0"
          v-bind:existing-reactions="existingReactions"
          v-bind:myUsername="commentsResponse.myUsername"
          @update:newReaction="createNewReaction"
          @update:deleteReaction="deleteExistingReaction"
          class="file-level-reactions"
        >
        </reaction-element>
        <h3 class="social-margin-left">Comments</h3>
        <div class="comments">
          <comment-element
            v-for="comment in commentsResponse.comments"
            :key="comment.id"
            v-bind:id="comment.id"
            v-bind:filePath="comment.filePath"
            v-bind:commentText="comment.commentText"
            v-bind:userName="comment.userName"
            v-bind:createdTime="comment.createdTime"
            v-bind:myUsername="commentsResponse.myUsername"
            class="comment"
            @update:deleteComment="deleteComment"
            @update:newReaction="createNewReaction"
            v-bind:reactionsAvailable="reactionsAvailable"
            v-bind:existing-reactions="existingReactions"
            @update:deleteReaction="deleteExistingReaction"
          >
          </comment-element>
        </div>
        <div class="commentEntering">
          <div class="commentEnteringHeader">New comment...</div>
          <textarea
            class="input input--block"
            v-model="commentInput"
            id="commentEnteringTextBox"
            v-on:focus="commentEnteringOnFocus"
            v-on:blur="commentEnteringOnBlur"
          />
          <div class="commentEnteringButtonDiv">
            <button class="button button--flat" @click="postComment">
              Post
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { reactions as reactionsApi } from "@/api";
import { comments as commentsApi } from "@/api";
import { files as filesApi } from "@/api";
import moment from "moment";
import CommentElement from "@/components/files/Comment.vue";
import ReactionElement from "@/components/files/Reaction.vue";

export default {
  name: "social",
  components: {
    CommentElement,
    ReactionElement,
  },
  props: [
    "filePath",
    "showComments",
    "screenWidth",
  ],
  async mounted() {
    const bodyElem = document.body;
    bodyElem.classList.add("noscroll");
    this.updateSocial();
    this.reactionsAvailable = await reactionsApi.fetchAvailableReactions();
    this.fetchAllReactionsForCurrentFile();
  },
  beforeDestroy() {
    const bodyElem = document.body;
    bodyElem.classList.remove("noscroll");
  },
  data: function () {
    return {
      commentsResponse: {
        comments: [],
        myUsername: "",
      },
      commentsLoaded: false,
      commentInput: "",
      reactionsAvailable: [],
      existingReactions: [],
    };
  },
  computed: {
    isFullDesktopMode() {
      return this.screenWidth >= 1400;
    },
  },
  methods: {
    async updateSocial() {
      if (!this.commentsLoaded) {
        try {
          this.commentsResponse = await commentsApi.fetchComments(
            this.filePath
          );
          if (!this.commentsResponse.comments) {
            this.commentsResponse.comments = [];
          }
          this.commentsLoaded = true;
        } catch (e) {
          this.$showError(e);
        }
      }
    },
    formatFromNowTime: function (timeValue) {
      return moment(timeValue).fromNow();
    },
    formatIso8601AbsoluteTime: function (timeValue) {
      return moment(timeValue).format("YYYY-MM-DD HH:mm");
    },
    async postComment(event) {
      event.preventDefault();

      await commentsApi.createComment(this.commentInput, this.filePath);
      this.commentsLoaded = false;
      this.commentInput = "";
      this.updateSocial();
    },
    async deleteComment(deletingId) {
      await commentsApi.deleteComment(deletingId);
      this.commentsLoaded = false;
      this.commentInput = "";
      this.updateSocial();
    },
    async toggleComments(event, newState) {
      event.preventDefault();
      this.$emit("update:showComments", newState);
    },
    getCommentUserPictureImgSrc(comment) {
      const imagePath = encodeURI(
        `/filebrowser-social/users/${comment.userName}.png`
      );
      return filesApi.getDownloadURL({ path: imagePath }, "thumb");
    },
    async fetchAllReactionsForCurrentFile() {
      this.existingReactions = (
        await reactionsApi.fetchReactions(this.filePath)
      ).reactions;
    },
    async createNewReaction(newReaction) {
      await reactionsApi.postCreateReaction(
        this.filePath,
        newReaction.reactionValue,
        newReaction.commentId
      );
      this.fetchAllReactionsForCurrentFile();
    },
    async deleteExistingReaction(deletingId) {
      await reactionsApi.deleteReaction(deletingId);
      this.fetchAllReactionsForCurrentFile();
    },
    async commentEnteringOnFocus() {
      this.$emit("update:commentEntering", true);
    },
    async commentEnteringOnBlur() {
      this.$emit("update:commentEntering", false);
    },
  },
};
</script>
<style>
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
  font-style: italic;
  font-size: 0.75em;
  padding-bottom: 0.5em;
}

#commentEnteringTextBox {
  line-height: 1.15;
  min-height: 10em;
  resize: vertical;
}

.socialSection {
  margin-bottom: 8em;
}

.inner-social-section {
  margin: 1em;
}

.social-margin-left {
  margin-left: 0.5em;
}

.file-level-reactions {
  margin-bottom: 1em;
}

.file-level-reactions .reaction-grid {
  background: var(--surfacePrimary);
}
</style>

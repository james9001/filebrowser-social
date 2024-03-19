import { fetchURL } from "./utils";

export async function fetchComments(url) {
  const res = await fetchURL(`/api/comments${encodeURIComponent(url)}`, {});

  if (res.status === 200) {
    let data = await res.json();
    return data;
  } else {
    throw new Error(res.status);
  }
}

export async function createComment(comment, filePath) {
  const res = await fetchURL(`/api/comments`, {
    method: "POST",
    body: JSON.stringify({
      data: {
        commentText: comment,
        filePath: filePath,
      },
      which: [],
      what: "comment",
    }),
  });

  if (res.status !== 200) {
    throw new Error(res.status);
  }
}

export async function deleteComment(commentId) {
  const res = await fetchURL(`/api/comments/${commentId}`, {
    method: "DELETE",
  });

  if (res.status !== 200) {
    throw new Error(res.status);
  }
}

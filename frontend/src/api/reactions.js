import { fetchURL } from "./utils";

export async function fetchAvailableReactions() {
  const res = await fetchURL(`/api/reactions/available`, {});

  if (res.status === 200) {
    let data = (await res.json()).reactionsAvailable;
    return data;
  } else {
    throw new Error(res.status);
  }
}

export async function fetchReactions(url) {
  const res = await fetchURL(`/api/reactions${encodeURIComponent(url)}`, {});

  if (res.status === 200) {
    let data = await res.json();
    return data;
  } else {
    throw new Error(res.status);
  }
}

export async function postCreateReaction(
  contextFilePath,
  reactionValue,
  commentId
) {
  const res = await fetchURL(`/api/reactions`, {
    method: "POST",
    body: JSON.stringify({
      data: {
        contextFilePath: contextFilePath,
        reactionValue: reactionValue,
        commentId: commentId,
      },
      which: [],
      what: "create_reaction",
    }),
  });

  if (res.status !== 200) {
    throw new Error(res.status);
  }
}

export async function deleteReaction(reactionId) {
  const res = await fetchURL(`/api/reactions/${reactionId}`, {
    method: "DELETE",
  });

  if (res.status !== 200) {
    throw new Error(res.status);
  }
}

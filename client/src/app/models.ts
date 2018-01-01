// see go code models


// TODO: hack a tool to generate ts definitions from go struct

interface Board {
  id: number;
  name: string;
  shortCode: string;
  description: string;
  isNsfw: boolean;
  maxAttachmentSize: number;
  allowedAttachmentTypes: string; // csv
}

interface AttachmentContainer {
  attachmentOriginalFilename: string;
  attachmentLocation: string;
  attachementTnLocation: string;
  attachmentSize: number;
}

interface Thread extends AttachmentContainer {
  id: number;
  boardID: number;
  description: string;
  subject: string;
  postedByID: number;
  createdAt: Date;
}

interface Post extends AttachmentContainer {
  id: number;
  threadID: number;
  content: string;
  postedAt: Date;
  isHidden: boolean;
  postedByID: number;
}

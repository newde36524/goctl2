package ast

import "github.com/newde36524/goctl2/pkg/parser/api/token"

// KVExpr is a key value expression.
type KVExpr struct {
	// Key is the key of the key value expression.
	Key *TokenNode
	// Value is the value of the key value expression.
	Value *TokenNode
}

func (i *KVExpr) HasHeadCommentGroup() bool {
	return i.Key.HasHeadCommentGroup()
}

func (i *KVExpr) HasLeadingCommentGroup() bool {
	return i.Value.HasLeadingCommentGroup()
}

func (i *KVExpr) CommentGroup() (head, leading CommentGroup) {
	return i.Key.HeadCommentGroup, i.Value.LeadingCommentGroup
}

func (i *KVExpr) Format(prefix ...string) string {
	w := NewBufferWriter()
	w.Write(withNode(i.Key, i.Value), withPrefix(prefix...), withInfix(Indent), withRawText())
	return w.String()
}

func (i *KVExpr) End() token.Position {
	return i.Value.End()
}

func (i *KVExpr) Pos() token.Position {
	return i.Key.Pos()
}

func (i *KVExpr) exprNode() {}

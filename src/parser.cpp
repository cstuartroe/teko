#include <stdio.h>
#include <fstream>
#include "tagger.cpp"
#include "nodes.cpp"

enum Precedence {NoPrec, Compare, BoolOp, AddSub, MultDiv, Exp, TopPrec};

Precedence getPrec(string infix) {
    if      (infix == "+")  { return AddSub; }
    else if (infix == "-")  { return AddSub; }
    else if (infix == "*")  { return MultDiv; }
    else if (infix == "/")  { return MultDiv; }
    else if (infix == "^")  { return Exp; }
    else if (infix == "%")  { return MultDiv; }
    else if (infix == "&")  { return BoolOp; }
    else if (infix == "|")  { return BoolOp; }
    else if (infix == "in") { return BoolOp; }
    else if (infix == "to") { return BoolOp; }
    else if (infix == "==") { return Compare; }
    else if (infix == "!=") { return Compare; }
    else if (infix == "<")  { return Compare; }
    else if (infix == ">")  { return Compare; }
    else if (infix == "<=") { return Compare; }
    else if (infix == ">=") { return Compare; }
    else if (infix == "<:") { return Compare; }
    else { throw runtime_error("Invalid infix: " + infix); }
}

struct TekoParser {
    Tag* first_tag = 0;
    Tag* curr_tag = 0;
    vector<string> lines;
    Statement *first_stmt = 0;
    Statement *curr_stmt = 0;

    TekoParser(string filename) {
        Tokenizer toker = Tokenizer();

        ifstream t;
        t.open(filename);
        string line;
        while(!t.eof()) {
            getline(t,line);
            lines.push_back(line);
            toker.digest_line(line);
        }
        t.close();

        Token *curr_token = toker.start;

        while (curr_token->next != 0) {
            Tag *new_tag = from_token(*curr_token);
            if (curr_tag != 0) {
                curr_tag->next = new_tag;
            } else {
                first_tag = new_tag;
            }
            curr_tag = new_tag;
            curr_token = curr_token->next;
        }
    }

    void TekoParserException(string message, Tag t) {
        printf("%s\n", string(12,'-').c_str());
        printf("Teko Parsing Error: %s\n", message.c_str());
        printf("at line %d, column %d\n", t.line_number, t.col);
        printf("    %s\n", lines[t.line_number-1].c_str());
        printf("%s\n", (string(t.col+3,' ') + string((int) t.s.size(), '~')).c_str());
        exit (EXIT_FAILURE);
    }

    void printout() {
        curr_tag = first_tag;
        while (curr_tag != 0) {
            printf("%s\n", curr_tag->to_str().c_str());
            curr_tag = curr_tag->next;
        }
    }

    void parse() {
        curr_tag = first_tag;
        while (curr_tag != 0) {
            Statement *new_stmt = grab_statement();
            if (first_stmt == 0) {
                first_stmt = new_stmt;
            } else {
                curr_stmt->next = new_stmt;
            }
            curr_stmt = new_stmt;
            printf("%s", curr_stmt->to_str(0).c_str());
        }
    }

    Tag *advance() {
        Tag *prev = curr_tag;
        if (curr_tag != 0) {
            curr_tag = curr_tag->next;
        }
        return prev;
    }

    Tag *get_curr_tag() {
        if (curr_tag == 0) {
            Tag fake_tag = Tag();
            fake_tag.s = " ";
            fake_tag.line_number = lines.size();
            fake_tag.col = lines[lines.size()-1].size();
            TekoParserException("Unexpected EOF", fake_tag);
        } else {
            return curr_tag;
        }
    }

    Statement *grab_statement(bool expect_semicolon = true) {
        Statement *stmt;

        if (get_curr_tag()->type == AnnotationTag || curr_tag->type == VartypeTag
            || curr_tag->type == LetTag) {
            stmt = grab_declstmt();
        } else {
            ExpressionNode *expr1 = grab_expression(NoPrec);
            if (get_curr_tag()->type == LabelTag) {
                DeclarationStmt *decl_stmt = new DeclarationStmt();
                decl_stmt->first_tag = expr1->first_tag;
                decl_stmt->tekotype = expr1;
                decl_stmt->declist = grab_declarations();
                stmt = decl_stmt;
            } else {
                ExpressionStmt *expr_stmt = new ExpressionStmt();
                expr_stmt->first_tag = expr1->first_tag;
                expr_stmt->body = expr1;
                stmt = expr_stmt;
            }
        }

        if (expect_semicolon) {
            if (get_curr_tag()->type == SemicolonTag) {
                advance();
            } else {
                TekoParserException("Expected semicolon", *get_curr_tag());
            }
        }

        return stmt;
    }

    DeclarationStmt *grab_declstmt() {
        DeclarationStmt *decl_stmt = new DeclarationStmt();
        decl_stmt->first_tag = get_curr_tag();

        while (get_curr_tag()->type == AnnotationTag) {
            char which_annot = *curr_tag->val;
            if (decl_stmt->annots[which_annot] != 0) {
                TekoParserException("Duplicate annotation " + annotations[*curr_tag->val], *curr_tag);
            }
            advance();

            AnnotationNode *annot = new AnnotationNode();
            if (annotations_params[which_annot]) {
                bool cont = true;
                while (cont) {
                    annot->params.push_back(grab_expression(TopPrec));
                    if (get_curr_tag()->type == CommaTag) {
                        advance();
                    } else {
                        cont = false;
                    }
                }
            }
            decl_stmt->annots[which_annot] = annot;
        }

        bool still_needs_type = true;
        for (int i = 0; i < num_vartypes; i++) {
            if (get_curr_tag()->type == VartypeTag && *curr_tag->val == i) {
                still_needs_type = false;
                decl_stmt->vts[i] = true;
                advance();
            }
        }
        if (get_curr_tag()->type == VartypeTag) {
            TekoParserException("Variable typing out of order", *curr_tag);
        }

        if (get_curr_tag()->type == LetTag) {
            if (!still_needs_type) {
                TekoParserException("Don't use let with variable types", *curr_tag);
            } else {
                still_needs_type = false;
                advance();
            }
        }

        if (still_needs_type) {
            decl_stmt->tekotype = grab_expression(TopPrec);
        }

        if (get_curr_tag()->type == SetterTag) {
            TekoParserException("Did not declare type", *curr_tag);
        }

        decl_stmt->declist = grab_declarations();

        return decl_stmt;
    }

    vector<DeclarationNode*> grab_declarations() {
        vector<DeclarationNode*> declist;

        while (get_curr_tag()->type == LabelTag) {
            DeclarationNode *decl_expr = new DeclarationNode();
            decl_expr->first_tag = curr_tag;
            decl_expr->label = *((string*)curr_tag->val);
            advance();

            if (get_curr_tag()->type == OpenTag && *((Brace*)curr_tag->val) == paren) {
                decl_expr->argstruct = grab_struct();
            }

            if (get_curr_tag()->type == SetterTag) {
                string set = setters[*curr_tag->val];
                if (set == "=" || set == "=&") {
                    decl_expr->setter = *curr_tag->val;
                    advance();
                    decl_expr->value = grab_expression(NoPrec);
                } else if (set == "->") {
                    decl_expr->setter = *curr_tag->val;
                    advance();
                    decl_expr->value = grab_codeblock();
                } else {
                    TekoParserException("Declarations may not use updating setters", *curr_tag);
                }
            }

            if (get_curr_tag()->type == CommaTag) {
                advance();
            } else if (curr_tag->type != SemicolonTag) {
                TekoParserException("Invalid declaration", *curr_tag);
            }
            declist.push_back(decl_expr);
        }
        return declist;
    }

    ExpressionNode *grab_expression(Precedence prec) {
        ExpressionNode *expr;
        Tag *expr_first_tag = get_curr_tag();

        switch (curr_tag->type) {
        case LabelTag:  expr = new SimpleNode(curr_tag); advance(); break;
        case StringTag: expr = new SimpleNode(curr_tag); advance(); break;
        case CharTag:   expr = new SimpleNode(curr_tag); advance(); break;
        case IntTag:    expr = new SimpleNode(curr_tag); advance(); break;
        case RealTag:   expr = new SimpleNode(curr_tag); advance(); break;
        case BoolTag:   expr = new SimpleNode(curr_tag); advance(); break;
        case IfTag:     expr = grab_if(); break;
        case ForTag:    expr = grab_for(); break;
        case WhileTag:  expr = grab_while(); break;

        case OpenTag:   {
            if (*((Brace*)curr_tag->val) == paren) {
                expr = grab_struct();
                if (expr == 0) {
                    expr = grab_args();
                }
            } else {
                expr = grab_sequence();
            }
        } break;

        case PrefixTag: {
            PrefixNode *prefix_expr = new PrefixNode();
            prefix_expr->first_tag = curr_tag;
            prefix_expr->prefix = *curr_tag->val;
            prefix_expr->right = grab_expression(TopPrec);
            expr = prefix_expr;
        } break;
        default: TekoParserException("Illegal start to expression: " + curr_tag->s, *curr_tag);
        }
        expr->first_tag = expr_first_tag;
        expr = continue_expression(prec, expr);
        return expr;
    }

    ExpressionNode *continue_expression(Precedence prec, ExpressionNode *left_expr) {
        if (curr_tag == 0) {
            return left_expr;
        }

        if (curr_tag->type == AttrTag) {
            if (curr_tag->next != 0 && curr_tag->next->type == LabelTag) {
                AttrNode *attr_expr = new AttrNode();
                attr_expr->first_tag = curr_tag;
                attr_expr->left = left_expr;
                advance();
                attr_expr->label = curr_tag->s;
                advance();
                return continue_expression(prec, attr_expr);
            } else {
                curr_tag->type = SuffixTag;
                char *suf = new char(string_index(".", suffixes, num_suffixes));
                curr_tag->val = suf;
            }
        }

        switch (curr_tag->type) {
        case InfixTag: {
            Precedence new_prec = getPrec(infixes[*curr_tag->val]);
            if (new_prec <= prec) {
                return left_expr;
            } else {
                InfixNode *infix_expr = new InfixNode();
                infix_expr->first_tag = curr_tag;
                infix_expr->left = left_expr;
                infix_expr->infix = *curr_tag->val;
                advance();
                infix_expr->right = grab_expression(new_prec);
                return continue_expression(prec, infix_expr);
            }
        } break;

        case SetterTag: {
            if (prec > NoPrec) {
                return left_expr;
            } else {
                AssignmentNode *asst_expr = new AssignmentNode();
                asst_expr->first_tag = left_expr->first_tag;
                asst_expr->left = left_expr;
                asst_expr->setter = *curr_tag->val;
                advance();
                string set = setters[asst_expr->setter];
                if (set == "->") {
                    asst_expr->right = grab_codeblock();
                } else {
                    asst_expr->right = grab_expression(NoPrec);
                }
                return asst_expr;
            }
        } break;

        case OpenTag: {
            Brace b = *((Brace *) curr_tag->val);
            switch (b) {
            case paren: {
                CallNode *call_expr = new CallNode();
                call_expr->first_tag = left_expr->first_tag;
                call_expr->left = left_expr;
                call_expr->args = grab_args();
                return continue_expression(prec, call_expr);
            } break;

            default: {
                SliceNode *slice_expr = new SliceNode();
                slice_expr->first_tag = curr_tag;
                slice_expr->left = left_expr;
                slice_expr->slice = grab_sequence();
                return continue_expression(prec, slice_expr);
            } break;
            }
        } break;

        case SuffixTag: {
            SuffixNode *suff_expr = new SuffixNode();
            suff_expr->first_tag = curr_tag;
            suff_expr->left = left_expr;
            suff_expr->suffix = *curr_tag->val;
            advance();
            return continue_expression(prec, suff_expr);
        }

        default: return left_expr; break;
        }
    }

    SeqNode *grab_sequence(bool forensic = false) {
        SeqNode *seq_expr = new SeqNode();
        if (get_curr_tag()->type != OpenTag) { throw runtime_error("erroneously called grab_sequence"); }
        seq_expr->first_tag = curr_tag;
        seq_expr->brace = *((Brace*) curr_tag->val);
        advance();

        while (get_curr_tag()->type != CloseTag) {
            ExpressionNode *elem = grab_expression(NoPrec);
            if (get_curr_tag()->type == ColonTag) {
                MapNode *map_expr = new MapNode();
                map_expr->first_tag = elem->first_tag;
                map_expr->key = elem;
                advance();
                map_expr->value = grab_expression(NoPrec);
                elem = map_expr;
            }
            seq_expr->elems.push_back(elem);

            if (get_curr_tag()->type == CommaTag) {
                advance();
            } else if (get_curr_tag()->type != CloseTag) {
                if (forensic) { return 0; }                 // for use when grabbing a codeblock
                TekoParserException("Expected a comma", *get_curr_tag());
            }
        }

        if (*((Brace*) curr_tag->val) != seq_expr->brace) {
            TekoParserException("Mismatched brace", *curr_tag);
        } else {
            advance();
        }

        return seq_expr;
    }

    ArgsNode *grab_args() {
        ArgsNode *args_expr = new ArgsNode();
        if (get_curr_tag()->type != OpenTag || *((Brace*) curr_tag->val) != paren) { throw runtime_error("erroneously called grab_args"); }
        args_expr->first_tag = curr_tag;
        advance();

        while (get_curr_tag()->type != CloseTag) {
            ArgNode *arg = new ArgNode();
            ExpressionNode* expr = grab_expression(NoPrec);
            if (expr->expr_type == SimpleExpr && ((SimpleNode*) expr)->data_type == LabelExpr && curr_tag->type == SetterTag && setters[*curr_tag->val] == "=") {
                arg->label = ((SimpleNode*) expr)->to_str(0);
                advance();
                expr = grab_expression(NoPrec);
            }
            arg->value = expr;
            args_expr->args.push_back(arg);

            if (get_curr_tag()->type == CommaTag) {
                advance();
            } else if (get_curr_tag()->type != CloseTag) {
                TekoParserException("Expected a comma", *get_curr_tag());
            }
        }

        if (*((Brace*) curr_tag->val) != paren) {
            TekoParserException("Mismatched brace", *curr_tag);
        } else {
            advance();
        }

        return args_expr;
    }

    StructNode *grab_struct() {
        Tag* orig_tag = get_curr_tag();

        StructNode *struct_expr = new StructNode();
        if (curr_tag->type != OpenTag || *((Brace*) curr_tag->val) != paren) { throw runtime_error("erroneously called grab_struct"); }
        struct_expr->first_tag = curr_tag;
        advance();

        if (get_curr_tag()->type == CloseTag) {
            if (*((Brace*) curr_tag->val) != paren) {
                TekoParserException("Mismatched brace", *curr_tag);
            } else {
                advance();
                return struct_expr;
            }
        }

        ExpressionNode* expr = grab_expression(NoPrec);

        if (get_curr_tag()->type != LabelTag) {
            curr_tag = orig_tag;
            return 0;
        } else {
            curr_tag = orig_tag;
            advance();
        }

        while (get_curr_tag()->type != CloseTag) {
            StructElemNode *elem = new StructElemNode();
            elem->tekotype = grab_expression(NoPrec);
            if (curr_tag->type != LabelTag) {
                TekoParserException("Expected label", *curr_tag);
            } else {
                elem->label = curr_tag->s;
                advance();
            }

            if (get_curr_tag()->type == QMarkTag) {
                advance();
                elem->deflt = grab_expression(NoPrec);
            }

            struct_expr->elems.push_back(elem);

            if (get_curr_tag()->type == CommaTag) {
                advance();
            } else if (get_curr_tag()->type != CloseTag) {
                TekoParserException("Expected a comma", *get_curr_tag());
            }
        }

        if (*((Brace*) curr_tag->val) != paren) {
            TekoParserException("Mismatched brace", *curr_tag);
        } else {
            advance();
        }

        return struct_expr;
    }

    Codeblock *grab_codeblock() {
        Codeblock *cb = new Codeblock();
        cb->first_tag = get_curr_tag();
        if (curr_tag->type == OpenTag && *((Brace*)curr_tag->val) == curly) {
            Tag *orig_tag = curr_tag;
            ExpressionNode *possible_set = grab_sequence(true);
            if (possible_set == 0) {
                curr_tag = orig_tag;
                advance();

                while (get_curr_tag()->type != CloseTag) {
                    Statement *stmt = grab_statement();
                    cb->stmts.push_back(stmt);
                }

                if (*((Brace*) curr_tag->val) != curly) {
                    TekoParserException("Mismatched brace", *curr_tag);
                } else {
                    advance();
                }
            }

            else {
                curr_tag = orig_tag;
                Statement *stmt = grab_statement(false);
                cb->stmts.push_back(stmt);
            }
        }

        else {
            Statement *stmt = grab_statement(false);
            cb->stmts.push_back(stmt);
        }
        return cb;
    }

    IfNode *grab_if() {
        if (get_curr_tag()->type != IfTag) {
            throw runtime_error("Erroneously called grab_if");
        }

        IfNode* if_expr = new IfNode();
        if_expr->first_tag = curr_tag;
        advance();

        if_expr->condition = grab_expression(NoPrec);
        if_expr->then_block = grab_codeblock();
        if_expr->else_block = grab_else();
        return if_expr;
    }

    Codeblock *grab_else() {
        if (curr_tag == 0) {
            return 0;
        }

        if (curr_tag->type == ElseTag || curr_tag->type == ColonTag) {
            advance();
            return grab_codeblock();
        } else {
            return 0;
        }
    }

    ForNode *grab_for() {
        if (get_curr_tag()->type != ForTag) {
            throw runtime_error("Erroneously called grab_for");
        }

        ForNode* for_expr = new ForNode();
        for_expr->first_tag = curr_tag;
        advance();

        bool populated = false;
        Tag *orig_tag = get_curr_tag();

        if (curr_tag->type == OpenTag && *((Brace*)curr_tag->val) == paren) {
            advance();
            populated = populate_for(for_expr);
            if (get_curr_tag()->type == CloseTag && *((Brace*)curr_tag->val) == paren) {
                advance();
            } else {
                TekoParserException("Mismatched parentheses", *curr_tag);
            }
        }

        if (!populated) {
            curr_tag = orig_tag;
            populated = populate_for(for_expr);
        }

        if (!populated) {
            TekoParserException("Invalid for loop", *get_curr_tag());
        }

        for_expr->codeblock = grab_codeblock();
        return for_expr;
    }

    bool populate_for(ForNode* for_expr) {
        ExpressionNode *expr1 = grab_expression(TopPrec);
        if (get_curr_tag()->type == LabelTag) {
            for_expr->tekotype = expr1;
            for_expr->label = *((string*)curr_tag->val);
            advance();
        } else if (expr1->expr_type == SimpleExpr && ((SimpleNode*)expr1)->data_type == LabelExpr) {
            for_expr->label = expr1->to_str(0);
        } else {
            return false;
        }

        if (get_curr_tag()->type != InfixTag || infixes[*curr_tag->val] != "in") {
            return false;
        } else {
            advance();
        }

        for_expr->iterator = grab_expression(NoPrec);
        return true;
    }

    WhileNode *grab_while() {
        if (get_curr_tag()->type != WhileTag) {
            throw runtime_error("Erroneously called grab_while");
        }

        WhileNode *while_expr = new WhileNode();
        while_expr->first_tag = curr_tag;
        advance();

        while_expr->condition = grab_expression(NoPrec);
        while_expr->codeblock = grab_codeblock();
        return while_expr;
    }
};

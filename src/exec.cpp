#include "types.cpp"

struct Interpreter {
    TekoParser *parser;
    TekoTypeChecker ttc;
    TekoObject module;

    Interpreter(TekoParser *_parser) {
        printf("o\n");
        parser = _parser;
        ttc = TekoTypeChecker();
        module = TekoObject(TekoModule);
        module.name = "<main>";

        module.declare("obj",    TekoTypeType, false);
        module.declare("type",   TekoTypeType, false);
        module.declare("bool",   TekoTypeType, false);
        module.declare("int",    TekoTypeType, false);
        module.declare("real",   TekoTypeType, false);
        module.declare("char",   TekoTypeType, false);
        module.declare("module", TekoTypeType, false);

        module.set("obj",    TekoObjectType);
        module.set("type",   TekoTypeType);
        module.set("bool",   TekoBool);
        module.set("int",    TekoInt);
        module.set("real",   TekoReal);
        module.set("char",   TekoChar);
        module.set("module", TekoModule);
    }

    void TekoTypeException(string message, Node *n) {
        Tag t = *n->first_tag;
        printf("%s\n", string(12,'-').c_str());
        printf("Type Error: %s\n", message.c_str());
        printf("at line %d, column %d\n", t.line_number, t.col);
        printf("    %s\n", parser->lines[t.line_number-1].c_str());
        printf("%s\n", (string(t.col+3,' ') + string((int) t.s.size(), '~')).c_str());
        exit (EXIT_FAILURE);
    }

    TekoObject *evaluate(ExpressionNode *expr) {
        switch (expr->expr_type) {
            case SimpleExpr:     return evaluate_simple((SimpleNode*) expr);
            case CallExpr:       return evaluate_call((CallNode*) expr);
            case InfixExpr:      return evaluate_infix((InfixNode*) expr);
            case AssignmentExpr: return evaluate_assignment((AssignmentNode*) expr);
            default: printf("evaluation not yet implemented! %d\n", expr->expr_type);
        }
    }

    TekoObject *evaluate_simple(SimpleNode *expr) {
        switch (expr->data_type) {
            case LabelExpr:
                return module.get(*((string*) expr->val));
            case StringExpr: break;
            case CharExpr: break;
            case RealExpr: break;
            case IntExpr:
                {
                    TekoObject *n = new TekoObject(TekoInt);
                    n->val = expr->val;
                    return n;
                }
            case BoolExpr: break;
            case BitsExpr: break;
            case BytesExpr: break;
        }
    }

    TekoObject *evaluate_call(CallNode *expr) {
        return 0;
    }

    TekoObject *evaluate_infix(InfixNode *expr) {
        return 0;
    }

    TekoObject *evaluate_assignment(AssignmentNode *expr) {
        return 0;
    }

    void execute(Statement *stmt) {
        printf("%s\n", stmt->to_str(0).c_str());
        ttc.check(stmt);

        switch (stmt->stmt_type) {
            case DeclarationStmtType: execute_decl((DeclarationStmt*) stmt); break;
            case NamespaceStmtType:   execute_ns((NamespaceStmt*) stmt); break;
            case ExpressionStmtType:  execute_expr((ExpressionStmt*) stmt); break;
        }

        if (stmt->next != 0) {
            execute(stmt->next);
        }
    }

    void execute_decl(DeclarationStmt *stmt) {
        TekoObject *type_obj = evaluate(stmt->tekotype);
        if (!is_instance(type_obj, TekoTypeType)) {
            TekoTypeException("Does not evaluate to a type", stmt->tekotype);
        }
        TekoType *type = (TekoType*) type_obj;

        bool is_mutable = stmt->vts[1];

        for (int i = 0; i < stmt->declist.size(); i++) {
            DeclarationNode *decl = stmt->declist[i];

            module.declare(decl->label, type, is_mutable);

            if (decl->value != 0) {
                TekoObject *value = evaluate(decl->value);
                if (!is_instance(value, type)) {
                    TekoTypeException("expected " + type->to_str(), decl->value);
                }
            }            
        }
    }

    void execute_ns(NamespaceStmt *stmt) {

    }

    void execute_expr(ExpressionStmt *expr_stmt) {
        evaluate(expr_stmt->body);
    }
};

void TekoFunction::execute(Interpreter i8er, ArgsNode *args) {
    for (int i = 0; i < stmts.size(); i++) {
        i8er.execute(stmts[i]);
    }
}

void TekoPrintFunction::execute(Interpreter i8er, ArgsNode *args) {
    TekoObject *s = i8er.evaluate(args->args[0]->value);
    printf("%s\n", s->to_str().c_str());
}
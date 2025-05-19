
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context,session sqlx.Session, data *{{.upperStartCamelObject}}) (sql.Result,error) {
	data.DeleteTime = time.Unix(0,0)
	data.DelState = globalkey.DelStateNo
	{{if .withCache}}{{.keys}}
	sqlResult, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
	query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
	if session != nil{
		return session.ExecCtx(ctx,query,{{.expressionValues}})
	}
	return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}})
	if err != nil {
	    return nil, err
	}
	data.Id, err = sqlResult.LastInsertId()
	if err != nil {
        return nil, err
	}
	return sqlResult, nil
	{{else}}
	query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
	if session == nil{
		session = m.conn
	}
	sqlResult, err := session.ExecCtx(ctx, query, {{.expressionValues}})
	if err != nil {
	    return nil, err
	}
	data.Id, err = sqlResult.LastInsertId()
	if err != nil {
        return nil, err
	}
	return sqlResult, nil{{end}}
}


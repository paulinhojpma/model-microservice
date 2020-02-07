package database

// arquivo para armazenar as querys para o banco de dados
const (
	// SELECTS
	SQLGetEscolaByID = `select row_to_json(t) from (
		select id_escola as "idEscola", nome, cnpj,
			(
			select array_to_json(array_agg(row_to_json(d)))
			from (
				select id_unidade as "idUnidade", nome, (
					select row_to_json(en)
						from (
							select id_endereco as "idEndereco", logradouro, numero, complemento, uf, cidade, cep
							from public.endereco where id_endereco = un.id_endereco
							)en
				) as endereco
				from escola.unidade un
				where id_escola =e.id_escola AND ativo = true

			) d
			) as unidades
		from escola.escola e where id_escola = %(id_escola)d

		) t;`

	SQLGetEscolas = `select row_to_json(t)
										from (
										  select id_escola as "idEscola", nome, cnpj,
										    (
										      select array_to_json(array_agg(row_to_json(d)))
										      from (
										        select id_unidade as "idUnidade", nome, (
															select row_to_json(en)
																from (
																	select id_endereco as "idEndereco", logradouro, numero, complemento, uf, cidade, cep
																	from public.endereco where id_endereco = un.id_endereco
																	)en
														) as endereco
										        from escola.unidade un
										        where id_escola =e.id_escola AND ativo = true

										      ) d
										    ) as unidades
										  from escola.escola e

										) t;`
	SQLGetDisciplinas = `SELECT row_to_json(t)
							from (
								SELECT id_disciplina as "idDisciplina", nome,  descricao,
								(
									SELECT array_to_json(array_agg(row_to_json(d)))
									from (
										SELECT id_serie_disciplina as "idEmenta", ementa, carga_horaria as "cargaHoraria",
										(
											SELECT row_to_json(s)
											FROM(
													SELECT id_serie as "idSerie", nome, tipo
													FROM escola.serie
													where id_serie = sd.id_serie
												) s

										) as serie
										from escola.serie_disciplina sd
										where id_disciplina = disc.id_disciplina
									) d
								) as ementas
								from escola.disciplina disc
								where id_escola = %(id_escola)d
								)t;`

	SQLGetDisciplinaByID = `SELECT row_to_json(t)
														from (
															SELECT id_disciplina as "idDisciplina", nome,  descricao,
															(
																SELECT array_to_json(array_agg(row_to_json(d)))
																from (
																	SELECT id_serie_disciplina as "idEmenta", ementa, carga_horaria as "cargaHoraria",
																	(
																		SELECT row_to_json(s)
																		FROM(
																				SELECT id_serie as "idSerie", nome, tipo
																				FROM escola.serie
																				where id_serie = sd.id_serie
																			) s

																	) as serie
																	from escola.serie_disciplina sd
																	where id_disciplina = %(id_disciplina)d
																) d
															) as ementas
															from escola.disciplina disc
															where id_escola = %(id_escola)d
															)t;`
	//INSERTS

	SQLInsertEscola   = `INSERT INTO escola.escola( nome, cnpj) VALUES ( '%(nome)s', '%(cnpj)s') returning id_escola;`
	SQLInsertUnidade  = `INSERT INTO escola.unidade( nome, id_endereco, id_escola) VALUES ( '%(nome)s',  %(id_endereco)d, %(id_escola)d) returning id_unidade;`
	SQLInsertEndereco = `INSERT INTO public.endereco( logradouro, bairro, cidade, uf, cep, numero, complemento)
	VALUES ( '%(logradouro)s', '%(bairro)s', '%(cidade)s', '%(uf)s', '%(cep)s', '%(numero)s', '%(complemento)s') returning id_endereco;`
	SQLInsertDisciplina = `INSERT INTO escola.disciplina( id_escola, nome, descricao)	VALUES (%(id_escola)d, '%(nome)s', '%(descricao)s') returning id_disciplina;`
	SQLInsertEmenta     = `INSERT INTO escola.serie_disciplina(id_serie, id_disciplina, carga_horaria, ementa)	
							VALUES ( %(id_serie)d, %(id_disciplina)d, %(carga_horaria)d, '%(ementa)s') 
							ON CONFLICT ON CONSTRAINT unique_serie_disciplina 
							DO UPDATE SET  carga_horaria =  %(carga_horaria)d, ementa = '%(ementa)s', ativo= %(ativo)t 
							returning id_serie_disciplina;`

	// UPDATES
	SQLUpdateDisciplina = `UPDATE escola.disciplina
	SET  nome= '%(nome)s', descricao= '%(descricao)s'
	WHERE id_disciplina = %(id_disciplina)d AND id_escola = %(id_escola)d  returning id_disciplina;`
)

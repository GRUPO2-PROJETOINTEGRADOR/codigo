--Criando tabela lojas
create table lojas (
	id varchar(20) primary key,
	nome varchar(100) not null,
	categoria varchar(100)
);

--Criando Tabelas Relatório Tiller

create table tiller (
id serial primary key, --O type serial incrementa 1 a cada adição
id_loja varchar(20) not null, --LUC da loja
data date not null, --Data da vistoria
nota int not null, -- Nota na data para a LUC relacionada
foreign key (id_loja) references lojas(id) ON DELETE CASCADE -- O id_loja deve ser exatamente como id em lojas
);

select * from lojas order by id;
select * from tiller order by id;

insert into lojas 
	(id, nome, categoria)
	VALUES
	('QS-03','CASA BAUDUCCO - QUIOSQUE','ALIMENTAÇÃO');

insert into tiller
(id_loja,data, nota)
VALUES
('QS-03','2026-06-18', 725);

select lojas.id, lojas.nome, tiller.id ,tiller.data, tiller.nota --Seleciona as colunas a ser exibidas para relacionar a visualização
FROM lojas, tiller order by tiller.data;

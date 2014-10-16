#!/usr/bin/perl

use lib 'extlib/lib/perl5';
use Mojolicious::Lite;

app->config(hypnotoad => {
		listen => ['http://*:8088', 'http://[::]:8088'],
		proxy => 1,
	});

get '/' => sub {
	my $c = shift;
	$c->render(text => $c->tx->remote_address . "\n");
};

app->start;

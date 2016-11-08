/*
  edit.js
  @flow

  The site editor page.
*/

var React = require('react');

import type { APIInstance } from '../api/api';
declare var api: APIInstance;

var Tree = require('../components/tree');
var Editor = require('../components/editor');
var Preview = require('../components/preview');
var Tab = require('../components/tab');

class Edit extends React.Component {
  state: {
    markdown: string,
    html: string,
    pages: Array<Page>,
    showEditor: boolean,
    showPreview: boolean,
    unsavedChanges: boolean,
    selectedPage: Page
  };

  converter: any;

  constructor(props: any) {
    super(props);

    this.state = {
      html: "",
      markdown: "",
      pages: [],
      showEditor: true,
      showPreview: true,
      selectedPage: {
        name: "index.md",
        markdown: "",
        html: ""
      },
      unsavedChanges: false
    }

    this.converter = new window.showdown.Converter();
  }

  componentDidMount() {
    // Load the list of pages for this site.
    api.pages().then((pages) => {
      this.setState({
        selectedPage: pages[0],
        pages
      });
    })

    api.addListener("ctrl+s", "editor", () => {
      this.save();
    });
  }

  save() {
    this.state.selectedPage.markdown = this.state.markdown;
    this.state.selectedPage.html = this.state.html;
    return api.editPage(this.state.selectedPage).then(() => {
      this.setState({unsavedChanges: false});
    });
  }

  textEdited(markdown: string) {
    var html = this.converter.makeHtml(markdown);
    this.setState({
      markdown,
      html,
      unsavedChanges:(this.state.selectedPage.markdown != markdown)
    });
  }

  clickTab(tab: 'editor' | 'preview') {
    if(tab == 'editor') {
      // Shouldn't be able to hide the editor if there isn't
      // anything else being displayed.
      if(this.state.showPreview || !this.state.showEditor)
        this.setState({showEditor: !this.state.showEditor})
    } else if (tab == 'preview') {
      this.setState({showPreview: !this.state.showPreview})
    }
  }

  clickPage(page: Page) {
    this.save().then(() => {
      this.setState({
        markdown: page.markdown,
        html: page.html,
        selectedPage: page
      });
    })
  }

  render() {
    return (
      <div style={styles.container}>
        <Tree clickPage={this.clickPage.bind(this)} pages={this.state.pages} />
        <div style={styles.editWindow}>
          <div style={styles.tabs}>
            <Tab onClick={this.clickTab.bind(this, 'editor')} selected={this.state.showEditor} indicator={this.state.unsavedChanges} name={this.state.selectedPage.name} />
            <Tab onClick={this.clickTab.bind(this, 'preview')} selected={this.state.showPreview} name="preview" />
          </div>
          <div style={styles.window}>
            <Editor visible={this.state.showEditor} initialText={this.state.selectedPage.markdown} onChange={this.textEdited.bind(this)} />
            {this.state.showPreview?(
              <Preview html={this.state.html} />
            ):[]}
          </div>
        </div>
      </div>
    );
  }
}

const styles = {
  container: {
    display: 'flex',
    flexDirection: 'row',
    flex: 1
  },
  editWindow: {
    display: 'flex',
    flexDirection: 'column',
    flex: 1
  },
  tabs: {
    display: 'flex',
    flexDirection: 'row',
    backgroundColor: '#716669'
  },
  window: {
    flex: 1,
    display: 'flex',
    flexDirection: 'row'
  }
};

module.exports = Edit;
